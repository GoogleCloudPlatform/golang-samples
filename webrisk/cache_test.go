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
	"reflect"
	"testing"
	"time"

	pb "github.com/GoogleCloudPlatform/golang-samples/webrisk/internal/webrisk_proto"
	pt "github.com/golang/protobuf/ptypes"
)

func TestCacheLookup(t *testing.T) {
	now := time.Unix(1451436338, 951473000)
	mockNow := func() time.Time { return now }

	type cacheLookup struct {
		h   hashPrefix
		tds map[ThreatType]bool
		r   cacheResult
	}
	vectors := []struct {
		gotCache  *cache // The cache to apply the Purge and Lookup on
		wantCache *cache // The cache expected after Purge
		lookups   []cacheLookup
	}{{
		gotCache: &cache{
			pttls: map[hashPrefix]map[ThreatType]time.Time{
				"AAAABBBBBBBBBBBBBBBBBBBBBBBBBBBB": {
					1: now.Add(DefaultUpdatePeriod),
				},
				"ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ": {
					2: now.Add(-time.Minute),
					1: now.Add(-DefaultUpdatePeriod),
				},
			},
			nttls: map[hashPrefix]time.Time{
				"AAAA": now.Add(DefaultUpdatePeriod),
				"BBBB": now.Add(-time.Minute),
			},
			now: mockNow,
		},
		wantCache: &cache{
			pttls: map[hashPrefix]map[ThreatType]time.Time{
				"AAAABBBBBBBBBBBBBBBBBBBBBBBBBBBB": {
					1: now.Add(DefaultUpdatePeriod),
				},
			},
			nttls: map[hashPrefix]time.Time{
				"AAAA": now.Add(DefaultUpdatePeriod),
			},
			now: mockNow,
		},
		lookups: []cacheLookup{{
			h:   "AAAABBBBBBBBBBBBBBBBBBBBBBBBBBBB",
			tds: map[ThreatType]bool{1: true},
			r:   positiveCacheHit,
		}, {
			h:   "AAAACDCDCDCDCDCDCDCDCDCDCDCDCDCD",
			tds: nil,
			r:   negativeCacheHit,
		}, {
			h:   "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ",
			tds: nil,
			r:   cacheMiss,
		}},
	}, {
		gotCache: &cache{
			pttls: map[hashPrefix]map[ThreatType]time.Time{
				"AAAABBBBBBBBBBBBBBBBBBBBBBBBBBBB": {
					1: now.Add(-DefaultUpdatePeriod),
				},
				"ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ": {
					2: now.Add(-time.Minute),
					1: now.Add(-DefaultUpdatePeriod),
				},
			},
			nttls: map[hashPrefix]time.Time{
				"AAAA": now.Add(DefaultUpdatePeriod * 2),
				"BBBB": now.Add(-time.Minute),
			},
			now: mockNow,
		},
		wantCache: &cache{
			pttls: map[hashPrefix]map[ThreatType]time.Time{
				"AAAABBBBBBBBBBBBBBBBBBBBBBBBBBBB": {
					1: now.Add(-DefaultUpdatePeriod),
				},
			},
			nttls: map[hashPrefix]time.Time{
				"AAAA": now.Add(DefaultUpdatePeriod * 2),
			},
			now: mockNow,
		},
		lookups: []cacheLookup{{
			h:   "AAAABBBBBBBBBBBBBBBBBBBBBBBBBBBB",
			tds: nil,
			r:   cacheMiss,
		}, {
			h:   "AAAACDCDCDCDCDCDCDCDCDCDCDCDCDCD",
			tds: nil,
			r:   negativeCacheHit,
		}, {
			h:   "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ",
			tds: nil,
			r:   cacheMiss,
		}},
	}, {
		gotCache:  &cache{now: mockNow},
		wantCache: &cache{now: mockNow},
		lookups: []cacheLookup{{
			h:   "AAAABBBBBBBBBBBBBBBBBBBBBBBBBBBB",
			tds: nil,
			r:   cacheMiss,
		}, {
			h:   "AAAACDCDCDCDCDCDCDCDCDCDCDCDCDCD",
			tds: nil,
			r:   cacheMiss,
		}},
	}}

	for i, v := range vectors {
		for j, l := range v.lookups {
			gotTDs, gotR := v.gotCache.Lookup(l.h)
			if !reflect.DeepEqual(gotTDs, l.tds) {
				t.Errorf("test %d, lookup %d, threats mismatch:\ngot  %+v\nwant %+v", i, j, gotTDs, l.tds)
			}
			if gotR != l.r {
				t.Errorf("test %d, lookup %d, result mismatch: got %d, want %d", i, j, gotR, l.r)

			}
		}
		v.gotCache.Purge()
		if !reflect.DeepEqual(v.wantCache.pttls, v.gotCache.pttls) {
			t.Errorf("purge test %d, mismatching cache contents: PTTLS\ngot  %+v\nwant %+v", i, v.gotCache.pttls, v.wantCache.pttls)
		}
		if !reflect.DeepEqual(v.wantCache.nttls, v.gotCache.nttls) {
			t.Errorf("purge test %d, mismatching cache contents: NTTLS\ngot  %+v\nwant %+v", i, v.gotCache.nttls, v.wantCache.nttls)
		}
		for j, l := range v.lookups {
			gotTDs, gotR := v.gotCache.Lookup(l.h)
			if !reflect.DeepEqual(gotTDs, l.tds) {
				t.Errorf("purge test %d, lookup %d, threats mismatch:\ngot  %+v\nwant %+v", i, j, gotTDs, l.tds)
			}
			if gotR != l.r {
				t.Errorf("purge test %d, lookup %d, result mismatch: got %d, want %d", i, j, gotR, l.r)
			}
		}
	}
}

func TestCacheUpdate(t *testing.T) {
	now := time.Unix(1451436338, 951473000)
	mockNow := func() time.Time { return now }
	ts, _ := pt.TimestampProto(now.Add(1000 * time.Second))
	tft, _ := pt.Timestamp(ts)

	vectors := []struct {
		req       *pb.SearchHashesRequest
		resp      *pb.SearchHashesResponse
		gotCache  *cache
		wantCache *cache
	}{{
		req:  &pb.SearchHashesRequest{},
		resp: &pb.SearchHashesResponse{},
		gotCache: &cache{
			now: mockNow,
		},
		wantCache: &cache{pttls: map[hashPrefix]map[ThreatType]time.Time{},
			nttls: map[hashPrefix]time.Time{},
			now:   mockNow,
		},
	}, {
		req: &pb.SearchHashesRequest{
			ThreatTypes: []pb.ThreatType{0, 1, 2},
			HashPrefix:  []byte("aaaa"),
		},
		resp: &pb.SearchHashesResponse{
			Threats: []*pb.SearchHashesResponse_ThreatHash{{
				ThreatTypes: []pb.ThreatType{0, 1, 2},
				Hash:        []byte("aaaabbbbccccddddeeeeffffgggghhhh"),
				ExpireTime:  ts,
			}},
			NegativeExpireTime: ts,
		},
		gotCache: &cache{
			now: mockNow,
		},
		wantCache: &cache{
			pttls: map[hashPrefix]map[ThreatType]time.Time{
				"aaaabbbbccccddddeeeeffffgggghhhh": {
					0: tft,
					1: tft,
					2: tft,
				},
			},
			nttls: map[hashPrefix]time.Time{
				"aaaa": tft,
			},
			now: mockNow,
		},
	}}

	for i, v := range vectors {
		err := v.gotCache.Update(v.req, v.resp)
		if err != nil {
			t.Fatalf("gotCache update returned unexpected error %v", err)
		}
		if !reflect.DeepEqual(v.wantCache.pttls, v.gotCache.pttls) {
			t.Errorf("test %d, mismatching cache contents: PTTLS\ngot  %+v\nwant %+v", i, v.gotCache.pttls, v.wantCache.pttls)
		}
		if !reflect.DeepEqual(v.wantCache.nttls, v.gotCache.nttls) {
			t.Errorf("test %d, mismatching cache contents: NTTLS\ngot  %+v\nwant %+v", i, v.gotCache.nttls, v.wantCache.nttls)
		}
	}
}
