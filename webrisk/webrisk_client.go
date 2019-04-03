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

// Package webrisk implements a client for the Web Risk API v4.
// API v4 emphasizes efficient usage of the network for bandwidth-constrained
// applications such as mobile devices. It achieves this by maintaining a small
// portion of the server state locally such that some queries can be answered
// immediately without any network requests. Thus, fewer API calls made, means
// less bandwidth is used.
//
// At a high-level, the implementation does the following:
//
//	            hash(query)
//	                 |
//	            _____V_____
//	           |           | No
//	           | Database  |-----+
//	           |___________|     |
//	                 |           |
//	                 | Maybe?    |
//	            _____V_____      |
//	       Yes |           | No  V
//	     +-----|   Cache   |---->+
//	     |     |___________|     |
//	     |           |           |
//	     |           | Maybe?    |
//	     |      _____V_____      |
//	     V Yes |           | No  V
//	     +<----|    API    |---->+
//	     |     |___________|     |
//	     V                       V
//	(Yes, unsafe)            (No, safe)
//
// Essentially the query is presented to three major components: The database,
// the cache, and the API. Each of these may satisfy the query immediately,
// or may say that it does not know and that the query should be satisfied by
// the next component. The goal of the database and cache is to satisfy as many
// queries as possible to avoid using the API.
//
// Starting with a user query, a hash of the query is performed to preserve
// privacy regarded the exact nature of the query. For example, if the query
// was for a URL, then this would be the SHA256 hash of the URL in question.
//
// Given a query hash, we first check the local database (which is periodically
// synced with the global Web Risk API servers). This database will either
// tell us that the query is definitely safe, or that it does not have
// enough information.
//
// If we are unsure about the query, we check the local cache, which can be used
// to satisfy queries immediately if the same query had been made recently.
// The cache will tell us that the query is either safe, unsafe, or unknown
// (because the it's not in the cache or the entry expired).
//
// If we are still unsure about the query, then we finally query the API server,
// which is guaranteed to return to us an authoritative answer, assuming no
// networking failures.
package webrisk

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"sync/atomic"
	"time"

	pb "github.com/GoogleCloudPlatform/golang-samples/webrisk/internal/webrisk_proto"
)

const (
	// DefaultServerURL is the default URL for the Web Risk API.
	DefaultServerURL = "webrisk.googleapis.com"

	// DefaultUpdatePeriod is the default period for how often WebriskClient will
	// reload its blacklist database.
	DefaultUpdatePeriod = 30 * time.Minute

	// DefaultID and DefaultVersion are the default client ID and Version
	// strings to send with every API call.
	DefaultID      = "GoWebriskClient"
	DefaultVersion = "1.0.0"

	// DefaultRequestTimeout is the default amount of time a single
	// api request can take.
	DefaultRequestTimeout = time.Minute
)

// Errors specific to this package.
var (
	errClosed = errors.New("webrisk: handler is closed")
	errStale  = errors.New("webrisk: threat list is stale")
)

// ThreatType is an enumeration type for threats classes. Examples of threat
// classes are malware, social engineering, etc.
type ThreatType uint16

func (tt ThreatType) String() string { return pb.ThreatType(tt).String() }

// List of ThreatType constants.
const (
	ThreatTypeUnspecified       = ThreatType(pb.ThreatType_THREAT_TYPE_UNSPECIFIED)
	ThreatTypeMalware           = ThreatType(pb.ThreatType_MALWARE)
	ThreatTypeSocialEngineering = ThreatType(pb.ThreatType_SOCIAL_ENGINEERING)
	ThreatTypeUnwantedSoftware  = ThreatType(pb.ThreatType_UNWANTED_SOFTWARE)
)

// DefaultThreatLists is the default list of threat lists that WebriskClient
// will maintain. Do not modify this variable.
var DefaultThreatLists = []ThreatType{
	ThreatTypeMalware,
	ThreatTypeSocialEngineering,
	ThreatTypeUnwantedSoftware,
}

// A URLThreat is a specialized ThreatType for the URL threat
// entry type.
type URLThreat struct {
	Pattern string
	ThreatType
}

// Config sets up the WebriskClient object.
type Config struct {
	// ServerURL is the URL for the Web Risk API server.
	// If empty, it defaults to DefaultServerURL.
	ServerURL string

	// ProxyURL is the URL of the proxy to use for all requests.
	// If empty, the underlying library uses $HTTP_PROXY environment variable.
	ProxyURL string

	// APIKey is the key used to authenticate with the Web Risk API
	// service. This field is required.
	APIKey string

	// ID and Version are client metadata associated with each API request to
	// identify the specific implementation of the client.
	// They are similar in usage to the "User-Agent" in an HTTP request.
	// If empty, these default to DefaultID and DefaultVersion, respectively.
	ID      string
	Version string

	// DBPath is a path to a persistent database file.
	// If empty, WebriskClient operates in a non-persistent manner.
	// This means that blacklist results will not be cached beyond the lifetime
	// of the WebriskClient object.
	DBPath string

	// UpdatePeriod determines how often we update the internal list database.
	// If zero value, it defaults to DefaultUpdatePeriod.
	UpdatePeriod time.Duration

	// ThreatLists determines which threat lists that WebriskClient should
	// subscribe to. The threats reported by LookupURLs will only be ones that
	// are specified by this list.
	// If empty, it defaults to DefaultThreatLists.
	ThreatLists []ThreatType

	// RequestTimeout determines the timeout value for the http client.
	RequestTimeout time.Duration

	// Logger is an io.Writer that allows WebriskClient to write debug information
	// intended for human consumption.
	// If empty, no logs will be written.
	Logger io.Writer

	// compressionTypes indicates how the threat entry sets can be compressed.
	compressionTypes []pb.CompressionType

	api api
	now func() time.Time
}

// setDefaults configures Config to have default parameters.
// It reports whether the current configuration is valid.
func (c *Config) setDefaults() bool {
	if c.ServerURL == "" {
		c.ServerURL = DefaultServerURL
	}
	if len(c.ThreatLists) == 0 {
		c.ThreatLists = DefaultThreatLists
	}
	if c.UpdatePeriod <= 0 {
		c.UpdatePeriod = DefaultUpdatePeriod
	}
	if c.RequestTimeout <= 0 {
		c.RequestTimeout = DefaultRequestTimeout
	}
	if c.compressionTypes == nil {
		c.compressionTypes = []pb.CompressionType{pb.CompressionType_RAW, pb.CompressionType_RICE}
	}
	return true
}

func (c Config) copy() Config {
	c2 := c
	c2.ThreatLists = append([]ThreatType(nil), c.ThreatLists...)
	c2.compressionTypes = append([]pb.CompressionType(nil), c.compressionTypes...)
	return c2
}

// WebriskClient is a client implementation of API v4.
//
// It provides a set of lookup methods that allows the user to query whether
// certain entries are considered a threat. The implementation manages all of
// local database and caching that would normally be needed to interact
// with the API server.
type WebriskClient struct {
	stats  Stats // Must be first for 64-bit alignment on non 64-bit systems.
	config Config
	api    api
	db     database
	c      cache

	lists map[ThreatType]bool

	log *log.Logger

	closed uint32
	done   chan bool // Signals that the updater routine should stop
}

// Stats records statistics regarding WebriskClient's operation.
type Stats struct {
	QueriesByDatabase int64         // Number of queries satisfied by the database alone
	QueriesByCache    int64         // Number of queries satisfied by the cache alone
	QueriesByAPI      int64         // Number of queries satisfied by an API call
	QueriesFail       int64         // Number of queries that could not be satisfied
	DatabaseUpdateLag time.Duration // Duration since last *missed* update. 0 if next update is in the future.
}

// NewWebriskClient creates a new WebriskClient.
//
// The conf struct allows the user to configure many aspects of the
// WebriskClient's operation.
func NewWebriskClient(conf Config) (*WebriskClient, error) {
	conf = conf.copy()
	if !conf.setDefaults() {
		return nil, errors.New("webrisk: invalid configuration")
	}

	// Create the SafeBrowsing object.
	if conf.api == nil {
		var err error
		conf.api, err = newNetAPI(conf.ServerURL, conf.APIKey, conf.ProxyURL)
		if err != nil {
			return nil, err
		}
	}
	if conf.now == nil {
		conf.now = time.Now
	}
	sb := &WebriskClient{
		config: conf,
		api:    conf.api,
		c:      cache{now: conf.now},
	}

	// TODO: Verify that config.ThreatLists is a subset of the list obtained
	// by "/v4/threatLists" API endpoint.

	// Convert threat lists slice to a map for O(1) lookup.
	sb.lists = make(map[ThreatType]bool)
	for _, td := range conf.ThreatLists {
		sb.lists[td] = true
	}

	// Setup the logger.
	w := conf.Logger
	if conf.Logger == nil {
		w = ioutil.Discard
	}
	sb.log = log.New(w, "webrisk: ", log.Ldate|log.Ltime|log.Lshortfile)

	delay := time.Duration(0)
	// If database file is provided, use that to initialize.
	if !sb.db.Init(&sb.config, sb.log) {
		ctx, cancel := context.WithTimeout(context.Background(), sb.config.RequestTimeout)
		delay, _ = sb.db.Update(ctx, sb.api)
		cancel()
	} else {
		if age := sb.db.SinceLastUpdate(); age < sb.config.UpdatePeriod {
			delay = sb.config.UpdatePeriod - age
		}
	}

	// Start the background list updater.
	sb.done = make(chan bool)
	go sb.updater(delay)
	return sb, nil
}

// Status reports the status of WebriskClient. It returns some statistics
// regarding the operation, and an error representing the status of its
// internal state. Most errors are transient and will recover themselves
// after some period.
func (sb *WebriskClient) Status() (Stats, error) {
	stats := Stats{
		QueriesByDatabase: atomic.LoadInt64(&sb.stats.QueriesByDatabase),
		QueriesByCache:    atomic.LoadInt64(&sb.stats.QueriesByCache),
		QueriesByAPI:      atomic.LoadInt64(&sb.stats.QueriesByAPI),
		QueriesFail:       atomic.LoadInt64(&sb.stats.QueriesFail),
		DatabaseUpdateLag: sb.db.UpdateLag(),
	}
	return stats, sb.db.Status()
}

// WaitUntilReady blocks until the database is not in an error state.
// Returns nil when the database is ready. Returns an error if the provided
// context is canceled or if the WebriskClient instance is Closed.
func (sb *WebriskClient) WaitUntilReady(ctx context.Context) error {
	if atomic.LoadUint32(&sb.closed) == 1 {
		return errClosed
	}
	select {
	case <-sb.db.Ready():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-sb.done:
		return errClosed
	}
}

// LookupURLs looks up the provided URLs. It returns a list of threats, one for
// every URL requested, and an error if any occurred. It is safe to call this
// method concurrently.
//
// The outer dimension is across all URLs requested, and will always have the
// same length as urls regardless of whether an error occurs or not.
// The inner dimension is across every fragment that a given URL produces.
// For some URL at index i, one can check for a hit on any blacklist by
// checking if len(threats[i]) > 0.
// The ThreatEntryType field in the inner ThreatType will be set to
// ThreatEntryType_URL as this is a URL lookup.
//
// If an error occurs, the caller should treat the threats list returned as a
// best-effort response to the query. The results may be stale or be partial.
func (sb *WebriskClient) LookupURLs(urls []string) (threats [][]URLThreat, err error) {
	threats, err = sb.LookupURLsContext(context.Background(), urls)
	return threats, err
}

// LookupURLsContext looks up the provided URLs. The request will be canceled
// if the provided Context is canceled, or if Config.RequestTimeout has
// elapsed. It is safe to call this method concurrently.
//
// See LookupURLs for details on the returned results.
func (sb *WebriskClient) LookupURLsContext(ctx context.Context, urls []string) (threats [][]URLThreat, err error) {
	ctx, cancel := context.WithTimeout(ctx, sb.config.RequestTimeout)
	defer cancel()

	threats = make([][]URLThreat, len(urls))

	if atomic.LoadUint32(&sb.closed) != 0 {
		return threats, errClosed
	}
	if err := sb.db.Status(); err != nil {
		sb.log.Printf("inconsistent database: %v", err)
		atomic.AddInt64(&sb.stats.QueriesFail, int64(len(urls)))
		return threats, err
	}

	hashes := make(map[hashPrefix]string)
	hash2idxs := make(map[hashPrefix][]int)

	// Construct the follow-up request being made to the server.
	// In the request, we only ask for partial hashes for privacy reasons.
	var reqs []*pb.SearchHashesRequest
	ttm := make(map[pb.ThreatType]bool)

	for i, url := range urls {
		urlhashes, err := generateHashes(url)
		if err != nil {
			sb.log.Printf("error generating urlhashes: %v", err)
			atomic.AddInt64(&sb.stats.QueriesFail, int64(len(urls)-i))
			return threats, err
		}

		for fullHash, pattern := range urlhashes {
			hash2idxs[fullHash] = append(hash2idxs[fullHash], i)
			_, alreadyRequested := hashes[fullHash]
			hashes[fullHash] = pattern

			// Lookup in database according to threat list.
			partialHash, unsureThreats := sb.db.Lookup(fullHash)
			if len(unsureThreats) == 0 {
				atomic.AddInt64(&sb.stats.QueriesByDatabase, 1)
				continue // There are definitely no threats for this full hash
			}

			// Lookup in cache according to recently seen values.
			cachedThreats, cr := sb.c.Lookup(fullHash)
			switch cr {
			case positiveCacheHit:
				// The cache remembers this full hash as a threat.
				// The threats we return to the client is the set intersection
				// of unsureThreats and cachedThreats.
				for _, td := range unsureThreats {
					if _, ok := cachedThreats[td]; ok {
						threats[i] = append(threats[i], URLThreat{
							Pattern:    pattern,
							ThreatType: td,
						})
					}
				}
				atomic.AddInt64(&sb.stats.QueriesByCache, 1)
			case negativeCacheHit:
				// This is cached as a non-threat.
				atomic.AddInt64(&sb.stats.QueriesByCache, 1)
				continue
			default:
				// The cache knows nothing about this full hash, so we must make
				// a request for it.
				if alreadyRequested {
					continue
				}
				for _, td := range unsureThreats {
					ttm[pb.ThreatType(td)] = true
				}

				tts := []pb.ThreatType{}
				for _, tt := range unsureThreats {
					tts = append(tts, pb.ThreatType(tt))
				}

				reqs = append(reqs, &pb.SearchHashesRequest{
					HashPrefix:  []byte(partialHash),
					ThreatTypes: tts,
				})
			}
		}
	}

	for _, req := range reqs {
		// Actually query the Web Risk API for exact full hash matches.
		resp, err := sb.api.HashLookup(ctx, req)
		if err != nil {
			sb.log.Printf("HashLookup failure: %v", err)
			atomic.AddInt64(&sb.stats.QueriesFail, 1)
			return threats, err
		}

		// Update the cache.
		sb.c.Update(req, resp)

		// Pull the information the client cares about out of the response.
		for _, threat := range resp.GetThreats() {
			fullHash := hashPrefix(threat.Hash)
			if !fullHash.IsFull() {
				continue
			}
			pattern, ok := hashes[fullHash]
			idxs, findidx := hash2idxs[fullHash]
			if findidx && ok {
				for _, td := range threat.ThreatTypes {
					if !sb.lists[ThreatType(td)] {
						continue
					}
					for _, idx := range idxs {
						threats[idx] = append(threats[idx], URLThreat{
							Pattern:    pattern,
							ThreatType: ThreatType(td),
						})
					}
				}
			}
		}
		atomic.AddInt64(&sb.stats.QueriesByAPI, 1)
	}
	return threats, nil
}

// TODO: Add other types of lookup when available.
//	func (sb *WebriskClient) LookupBinaries(digests []string) (threats []BinaryThreat, err error)
//	func (sb *WebriskClient) LookupAddresses(addrs []string) (threats [][]AddressThreat, err error)

// updater is a blocking method that periodically updates the local database.
// This should be run as a separate goroutine and will be automatically stopped
// when sb.Close is called.
func (sb *WebriskClient) updater(delay time.Duration) {
	for {
		sb.log.Printf("Next update in %v", delay)
		select {
		case <-time.After(delay):
			var ok bool
			ctx, cancel := context.WithTimeout(context.Background(), sb.config.RequestTimeout)
			if delay, ok = sb.db.Update(ctx, sb.api); ok {
				sb.log.Printf("background threat list updated")
				sb.c.Purge()
			}
			cancel()

		case <-sb.done:
			return
		}
	}
}

// Close cleans up all resources.
// This method must not be called concurrently with other lookup methods.
func (sb *WebriskClient) Close() error {
	if atomic.LoadUint32(&sb.closed) == 0 {
		atomic.StoreUint32(&sb.closed, 1)
		close(sb.done)
	}
	return nil
}
