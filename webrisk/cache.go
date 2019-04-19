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
	"fmt"
	"sync"
	"time"

	pb "github.com/GoogleCloudPlatform/golang-samples/webrisk/internal/webrisk_proto"
	pt "github.com/golang/protobuf/ptypes"
)

type cacheResult int

const (
	// positiveCacheHit indicates that the given hash matched an entry in the cache.
	// The caller must consider the match a threat and not contact the server.
	positiveCacheHit cacheResult = iota

	// negativeCacheHit indicates that the given hash did not match any entries
	// in the cache but its prefix matches the negative cache. The caller must
	// consider the given hash to be safe and not contact the server.
	negativeCacheHit

	// cacheMiss indicates that the given hash did not match any entry
	// in the cache. The caller should make a follow-up query to the server.
	cacheMiss

	// cacheError indicates that there was an error while looking up an
	// entry.
	cacheError
)

// cache caches results from API calls to SearchHashesRequest to reduce
// network calls for recently requested items. Since the global blacklist is
// constantly changing, the Web Risk API defines TTLs for how long entries
// can stay alive in the cache.
type cache struct {
	sync.RWMutex

	// pttls maps full hashes and a ThreatType to a positive time-to-live.
	// For a given full hash, the known threats are all ThreatTypes that
	// map to valid TTLs (i.e. in the future).
	pttls map[hashPrefix]map[ThreatType]time.Time

	// nttls maps partial hashes to a negative time-to-live.
	// If this is still valid (i.e. in the future), then this indicates that
	// there are *no* threats under the given partial hash, unless there exist
	// ThreatTypes with a valid positive TTL for that hash.
	nttls map[hashPrefix]time.Time

	now func() time.Time
}

// Update updates the cache according to the request that was made to the server
// and the response given back.
func (c *cache) Update(req *pb.SearchHashesRequest, resp *pb.SearchHashesResponse) error {
	c.Lock()
	defer c.Unlock()

	if c.pttls == nil {
		c.pttls = make(map[hashPrefix]map[ThreatType]time.Time)
		c.nttls = make(map[hashPrefix]time.Time)
	}

	// Insert each threat match into the cache by full hash.
	for _, threat := range resp.GetThreats() {
		fullHash := hashPrefix(threat.Hash)
		if !fullHash.IsFull() {
			continue
		}
		if c.pttls[fullHash] == nil {
			c.pttls[fullHash] = make(map[ThreatType]time.Time)
		}
		for _, tt := range threat.ThreatTypes {
			var err error
			c.pttls[fullHash][ThreatType(tt)], err = pt.Timestamp(threat.ExpireTime)
			if err != nil {
				return fmt.Errorf("pt.Timestamp: %v", err)
			}
		}
	}

	// Insert negative TTLs for partial hashes.
	if resp.GetNegativeExpireTime() != nil {
		nttl, _ := pt.Timestamp(resp.GetNegativeExpireTime())
		partialHash := hashPrefix(req.HashPrefix)
		c.nttls[partialHash] = nttl
	}
	return nil
}

// Lookup looks up a full hash and returns a set of ThreatTypes and the
// validity of the result.
func (c *cache) Lookup(hash hashPrefix) (map[ThreatType]bool, cacheResult) {
	if !hash.IsFull() {
		return nil, cacheError
	}

	c.Lock()
	defer c.Unlock()
	now := c.now()

	// Check all entries to see if there *is* a threat.
	threats := make(map[ThreatType]bool)
	threatTTLs := c.pttls[hash]
	for td, pttl := range threatTTLs {
		if pttl.After(now) {
			threats[td] = true
		} else {
			// The PTTL has expired, we should ask the server what's going on.
			return nil, cacheMiss
		}
	}
	if len(threats) > 0 {
		// So long as there are valid threats, we report them. The positive TTL
		// takes precedence over the negative TTL at the partial hash level.
		return threats, positiveCacheHit
	}

	// Check the negative TTLs to see if there are *no* threats.
	for i := minHashPrefixLength; i <= maxHashPrefixLength; i++ {
		if nttl, ok := c.nttls[hash[:i]]; ok {
			if nttl.After(now) {
				return nil, negativeCacheHit
			}
		}
	}

	// The cache has no information; it is a *possible* threat.
	return nil, cacheMiss
}

// Purge purges all expired entries from the cache.
func (c *cache) Purge() {
	c.Lock()
	defer c.Unlock()
	now := c.now()

	// Nuke all threat entries based on their positive TTL.
	for fullHash, threatTTLs := range c.pttls {
		for td, pttl := range threatTTLs {
			if now.After(pttl) {
				del := true
				for i := minHashPrefixLength; i <= maxHashPrefixLength; i++ {
					if nttl, ok := c.nttls[fullHash[:i]]; ok {
						if nttl.After(pttl) {
							del = false
							break
						}
					}
				}
				if del {
					delete(threatTTLs, td)
				}
			}
		}
		if len(threatTTLs) == 0 {
			delete(c.pttls, fullHash)
		}
	}

	// Nuke all partial hashes based on their negative TTL.
	for partialHash, nttl := range c.nttls {
		if now.After(nttl) {
			delete(c.nttls, partialHash)
		}
	}
}
