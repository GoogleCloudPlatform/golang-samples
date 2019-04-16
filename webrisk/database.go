// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package webrisk

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/gob"
	"errors"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	pb "github.com/GoogleCloudPlatform/golang-samples/webrisk/internal/webrisk_proto"
	pt "github.com/golang/protobuf/ptypes"
)

// jitter is the maximum amount of time that we expect an API list update to
// actually take. We add this time to the update period time to give some
// leeway before declaring the database as stale.
const (
	maxRetryDelay  = 24 * time.Hour
	baseRetryDelay = 15 * time.Minute
	jitter         = 30 * time.Second
)

// database tracks the state of the threat lists published by the Webrisk API.
// Since the global blacklist is constantly changing, the contents of the
// database needs to be periodically synced with the Webrisk servers in
// order to provide protection for the latest threats.
//
// The process for updating the database is as follows:
//	* At startup, if a database file is provided, then load it. If loaded
//	properly (not corrupted and not stale), then set tfu as the contents.
//	Otherwise, pull a new threat list from the Web Risk API.
//	* Periodically, synchronize the database with the Web Risk API.
//	This uses the Version Token fields to update only parts of the threat list that have
//	changed since the last sync.
//	* Anytime tfu is updated, generate a new tfl.
//
// The process for querying the database is as follows:
//	* Check if the requested full hash matches any partial hash in tfl.
//	If a match is found, return a set of ThreatTypes with a partial match.
type database struct {
	ml sync.RWMutex // Protects tfl, err, and last
	// threatsForLookup maps ThreatTypes to sets of partial hashes.
	// This data structure is in a format that is easily queried.
	tfl  threatsForLookup
	err  error     // Last error encountered
	last time.Time // Last time the threat list were synced

	config *Config
	// threatsForUpdate maps ThreatTypes to lists of partial hashes.
	// This data structure is in a format that is easily updated by the API.
	// It is also the form that is written to disk.
	tfu threatsForUpdate
	mu  sync.Mutex // Protects tfu

	readyCh         chan struct{} // Used for waiting until not in an error state.
	updateAPIErrors uint          // Number of times we attempted to contact the api and failed

	log *log.Logger
}

type threatsForUpdate map[ThreatType]partialHashes
type partialHashes struct {
	// Since the Hashes field is only needed when storing to disk and when
	// updating, this field is cleared except for when it is in use.
	// This is done to reduce memory usage as the contents of this can be
	// regenerated from the tfl.
	Hashes hashPrefixes

	SHA256 []byte // The SHA256 over Hashes
	State  []byte // Arbitrary binary blob to synchronize state with API
}

type threatsForLookup map[ThreatType]hashSet

// databaseFormat is a light struct used only for gob encoding and decoding.
// As written to disk, the format of the database file is basically the gzip
// compressed version of the gob encoding of databaseFormat.
type databaseFormat struct {
	Table threatsForUpdate
	Time  time.Time
}

// Init initializes the database from the specified file in config.DBPath.
// It reports true if the database was successfully loaded. If it reports false
// use Status for more details on the failure.
func (db *database) Init(config *Config, logger *log.Logger) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.setError(errors.New("not intialized"))
	db.config = config
	db.log = logger
	if db.config.DBPath == "" {
		db.log.Printf("no database file specified")
		db.setError(errors.New("no database loaded"))
		return false
	}
	dbf, err := loadDatabase(db.config.DBPath)
	if err != nil {
		db.log.Printf("load failure: %v", err)
		db.setError(err)
		return false
	}
	// Validate that the database threat list stored on disk is not too stale.
	if db.isStale(dbf.Time) {
		db.log.Printf("database loaded is stale")
		db.ml.Lock()
		defer db.ml.Unlock()
		db.setStale()
		return false
	}
	// Validate that the database threat list stored on disk is at least a
	// superset of the specified configuration.
	tfuNew := make(threatsForUpdate)
	for _, td := range db.config.ThreatLists {
		if row, ok := dbf.Table[td]; ok {
			tfuNew[td] = row
		} else {
			db.log.Printf("database configuration mismatch, missing %v", td)
			db.setError(errors.New("database configuration mismatch"))
			return false
		}
	}
	db.tfu = tfuNew
	db.generateThreatsForLookups(dbf.Time)
	return true
}

// Status reports the health of the database. The database is considered faulted
// if there was an error during update or if the last update has gone stale. If
// in a faulted state, the db may repair itself on the next Update.
func (db *database) Status() error {
	db.ml.RLock()
	defer db.ml.RUnlock()

	if db.err != nil {
		return db.err
	}
	if db.isStale(db.last) {
		db.setStale()
		return db.err
	}
	return nil
}

// UpdateLag reports the amount of time in between when we expected to run
// a database update and the current time
func (db *database) UpdateLag() time.Duration {
	lag := db.SinceLastUpdate()
	if lag < db.config.UpdatePeriod {
		return 0
	}
	return lag - db.config.UpdatePeriod
}

// SinceLastUpdate gives the duration since the last database update
func (db *database) SinceLastUpdate() time.Duration {
	db.ml.RLock()
	defer db.ml.RUnlock()

	return db.config.now().Sub(db.last)
}

// Ready returns a channel that's closed when the database is ready for queries.
func (db *database) Ready() <-chan struct{} {
	return db.readyCh
}

// Update synchronizes the local threat lists with those maintained by the
// global Web Risk API servers. If the update is successful, Status should
// report a nil error.
func (db *database) Update(ctx context.Context, api api) (time.Duration, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Construct and make the requests.
	var s []*pb.ComputeThreatListDiffRequest
	for _, td := range db.config.ThreatLists {
		var state []byte
		if row, ok := db.tfu[td]; ok {
			state = row.State
		}

		s = append(s, &pb.ComputeThreatListDiffRequest{
			ThreatType: pb.ThreatType(td),
			Constraints: &pb.ComputeThreatListDiffRequest_Constraints{
				SupportedCompressions: db.config.compressionTypes},
			VersionToken: state,
		})
	}

	var resps []*pb.ComputeThreatListDiffResponse

	// add jitter to wait time to avoid all servers lining up
	nextUpdateWait := db.config.UpdatePeriod + time.Duration(rand.Int31n(60)-30)*time.Second
	last := db.config.now()
	for _, req := range s {
		// Query the API for the threat list and update the database.
		resp, err := api.ListUpdate(ctx, req)
		if err != nil {
			db.log.Printf("ListUpdate failure (%d): %v", db.updateAPIErrors+1, err)
			db.setError(err)
			// backoff strategy: MIN((2**N-1 * 15 minutes) * (RAND + 1), 24 hours)
			n := 1 << db.updateAPIErrors
			delay := time.Duration(float64(n) * (rand.Float64() + 1) * float64(baseRetryDelay))
			if delay > maxRetryDelay {
				delay = maxRetryDelay
			}
			db.updateAPIErrors++
			return delay, false
		}
		resps = append(resps, resp)
		if resp.RecommendedNextDiff != nil {
			ndiff, _ := pt.Timestamp(resp.RecommendedNextDiff)
			serverMinWait := time.Duration(ndiff.Sub(time.Now()))
			if serverMinWait > nextUpdateWait {
				nextUpdateWait = serverMinWait
				db.log.Printf("Server requested next update in %v", nextUpdateWait)
			}
		}
	}

	// If for some reason we missed a request or didn't get a response the
	// rest of the logic may fail.
	if len(s) != len(resps) {
		db.setError(errors.New("mismatch between requests sent and responses received"))
		return nextUpdateWait, false
	}

	db.updateAPIErrors = 0
	// Update the threat database with the response.
	db.generateThreatsForUpdate()
	for i, resp := range resps {
		// Assume a 1:1 correspondence between request and response
		if err := db.tfu.update(resp, ThreatType(s[i].ThreatType)); err != nil {
			db.setError(err)
			db.log.Printf("update failure: %v", err)
			db.tfu = nil
			return nextUpdateWait, false
		}
	}

	dbf := databaseFormat{make(threatsForUpdate), last}
	for td, phs := range db.tfu {
		// Copy of partialHashes before generateThreatsForLookups clobbers it.
		dbf.Table[td] = phs
	}

	db.generateThreatsForLookups(last)

	// Regenerate the database and store it.
	if db.config.DBPath != "" {
		// Semantically, we ignore save errors, but we do log them.
		if err := saveDatabase(db.config.DBPath, dbf); err != nil {
			db.log.Printf("save failure: %v", err)
		}
	}
	return nextUpdateWait, true
}

// Lookup looks up the full hash in the threat list and returns a partial
// hash and a set of ThreatTypes that may match the full hash.
func (db *database) Lookup(hash hashPrefix) (h hashPrefix, tds []ThreatType) {
	if !hash.IsFull() {
		panic("hash is not full")
	}

	db.ml.RLock()
	for td, hs := range db.tfl {
		if n := hs.Lookup(hash); n > 0 {
			h = hash[:n]
			tds = append(tds, td)
		}
	}
	db.ml.RUnlock()
	return h, tds
}

// setError clears the database state and sets the last error to be err.
//
// This assumes that the db.mu lock is already held.
func (db *database) setError(err error) {
	db.tfu = nil

	db.ml.Lock()
	if db.err == nil {
		db.readyCh = make(chan struct{})
	}
	db.tfl, db.err, db.last = nil, err, time.Time{}
	db.ml.Unlock()
}

// isStale checks whether the last successful update should be considered stale.
// Staleness is defined as being older than two of the configured update periods
// plus jitter.
func (db *database) isStale(lastUpdate time.Time) bool {
	return db.config.now().Sub(lastUpdate) > 2*(db.config.UpdatePeriod+jitter)
}

// setStale sets the error state to a stale message, without clearing
// the database state.
//
// This assumes that the db.ml lock is already held.
func (db *database) setStale() {
	if db.err == nil {
		db.readyCh = make(chan struct{})
	}
	db.err = errStale
}

// clearError clears the db error state, and unblocks any callers of
// WaitUntilReady.
//
// This assumes that the db.mu lock is already held.
func (db *database) clearError() {
	db.ml.Lock()
	defer db.ml.Unlock()

	if db.err != nil {
		close(db.readyCh)
	}
	db.err = nil
}

// generateThreatsForUpdate regenerates the threatsForUpdate hashes from
// the threatsForLookup. We do this to avoid holding onto the hash lists for
// a long time, needlessly occupying lots of memory.
//
// This assumes that the db.mu lock is already held.
func (db *database) generateThreatsForUpdate() {
	if db.tfu == nil {
		db.tfu = make(threatsForUpdate)
	}

	db.ml.RLock()
	for td, hs := range db.tfl {
		phs := db.tfu[td]
		phs.Hashes = hs.Export()
		db.tfu[td] = phs
	}
	db.ml.RUnlock()
}

// generateThreatsForLookups regenerates the threatsForLookup data structure
// from the threatsForUpdate data structure and stores the last timestamp.
// Since the hashes are effectively stored as a set inside the threatsForLookup,
// we clear out the hashes slice in threatsForUpdate so that it can be GCed.
//
// This assumes that the db.mu lock is already held.
func (db *database) generateThreatsForLookups(last time.Time) {
	tfl := make(threatsForLookup)
	for td, phs := range db.tfu {
		var hs hashSet
		hs.Import(phs.Hashes)
		tfl[td] = hs

		phs.Hashes = nil // Clear hashes to keep memory usage low
		db.tfu[td] = phs
	}

	db.ml.Lock()
	wasBad := db.err != nil
	db.tfl, db.last = tfl, last
	db.ml.Unlock()

	if wasBad {
		db.clearError()
		db.log.Printf("database is now healthy")
	}
}

// saveDatabase saves the database threat list to a file.
func saveDatabase(path string, db databaseFormat) (err error) {
	var file *os.File
	file, err = os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); err == nil {
			err = cerr
		}
	}()

	gz, err := gzip.NewWriterLevel(file, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer func() {
		if zerr := gz.Close(); err == nil {
			err = zerr
		}
	}()

	encoder := gob.NewEncoder(gz)
	if err = encoder.Encode(db); err != nil {
		return err
	}
	return nil
}

// loadDatabase loads the database state from a file.
func loadDatabase(path string) (db databaseFormat, err error) {
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		return db, err
	}
	defer func() {
		if cerr := file.Close(); err == nil {
			err = cerr
		}
	}()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return db, err
	}
	defer func() {
		if zerr := gz.Close(); err == nil {
			err = zerr
		}
	}()

	decoder := gob.NewDecoder(gz)
	if err = decoder.Decode(&db); err != nil {
		return db, err
	}
	for _, dv := range db.Table {
		if !bytes.Equal(dv.SHA256, dv.Hashes.SHA256()) {
			return db, errors.New("webrisk: threat list SHA256 mismatch")
		}
	}
	return db, nil
}

// update updates the threat list according to the API response.
func (tfu threatsForUpdate) update(resp *pb.ComputeThreatListDiffResponse, td ThreatType) error {
	phs, ok := tfu[td]

	removalQuantity := 0
	if resp.ResponseType == pb.ComputeThreatListDiffResponse_RESET {
		phs = partialHashes{}
	}
	if resp.Removals != nil {
		if resp.Removals.RawIndices != nil {
			removalQuantity += len(resp.Removals.RawIndices.Indices)
		}
		if resp.Removals.RiceIndices != nil {
			if resp.Removals.RiceIndices.EntryCount == 0 {
				removalQuantity++
			} else {
				removalQuantity += int(resp.Removals.RiceIndices.EntryCount)
			}
		}
		switch resp.ResponseType {
		case pb.ComputeThreatListDiffResponse_DIFF:
			if !ok {
				return errors.New("webrisk: partial update received for non-existent key")
			}
		case pb.ComputeThreatListDiffResponse_RESET:
			if removalQuantity > 0 {
				return errors.New("webrisk: indices to be removed included in a full update")
			}
		default:
			return errors.New("webrisk: unknown response type")
		}

		// Hashes must be sorted for removal logic to work properly.
		phs.Hashes.Sort()

		idxs, err := decodeIndices(resp.Removals)
		if err != nil {
			return err
		}

		for _, i := range idxs {
			if i < 0 || i >= int32(len(phs.Hashes)) {
				return errors.New("webrisk: invalid removal index")
			}
			phs.Hashes[i] = ""
		}

		// If any removal was performed, compact the list of hashes.
		if removalQuantity > 0 {
			compactHashes := phs.Hashes[:0]
			for _, h := range phs.Hashes {
				if h != "" {
					compactHashes = append(compactHashes, h)
				}
			}
			phs.Hashes = compactHashes
		}
	}

	if resp.Additions != nil {

		hashes, err := decodeHashes(resp.Additions)
		if err != nil {
			return err
		}
		phs.Hashes = append(phs.Hashes, hashes...)
	}

	// Hashes must be sorted for SHA256 checksum to be correct.
	phs.Hashes.Sort()
	if err := phs.Hashes.Validate(); err != nil {
		return err
	}

	if cs := resp.GetChecksum(); cs != nil {
		phs.SHA256 = cs.Sha256
	}
	if !bytes.Equal(phs.SHA256, phs.Hashes.SHA256()) {
		return errors.New("webrisk: threat list SHA256 mismatch")
	}

	phs.State = resp.NewVersionToken
	tfu[td] = phs
	return nil
}
