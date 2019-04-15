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
	"context"
	"os"
	"testing"
	"time"

	pb "github.com/GoogleCloudPlatform/golang-samples/webrisk/internal/webrisk_proto"
)

// The system tests below are non-deterministic and operate by performing
// network requests against the Web Risk API servers. Thus, in order to
// operate they need the user's API key. This can be specified using the -apikey
// command-line flag when running the tests.
var apiKey = os.Getenv("WEBRISK_APIKEY")

func TestNetworkAPIUpdate(t *testing.T) {
	if apiKey == "" {
		t.Skip()
	}

	nm, err := newNetAPI(DefaultServerURL, apiKey, "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	req := &pb.ComputeThreatListDiffRequest{
		ThreatType: pb.ThreatType_MALWARE,
	}

	dat, err := nm.ListUpdate(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	hashes, err := decodeHashes(dat.GetAdditions())
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(hashes) && i < 10; i++ {
		hash := hashes[i]
		fullHashReq := &pb.SearchHashesRequest{
			ThreatTypes: []pb.ThreatType{pb.ThreatType_MALWARE},
			HashPrefix:  []byte(hash),
		}
		fullHashResp, err := nm.HashLookup(context.Background(), fullHashReq)
		if err != nil {
			t.Fatal(err)
		}
		if got := len(fullHashResp.GetThreats()); got < 1 {
			t.Fatalf("len(r.GetMatches()), got: %v, want: > 0", got)
		}
	}
}

func TestNetworkAPILookup(t *testing.T) {
	if apiKey == "" {
		t.Skip()
	}

	nm, err := newNetAPI(DefaultServerURL, apiKey, "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var c = pb.ComputeThreatListDiffRequest{
		ThreatType: pb.ThreatType_MALWARE,
	}
	url := "http://testwebrisk | w.appspot.com/apiv4/ANY_PLATFORM/MALWARE/URL/"
	hash := hashFromPattern(url)
	req := &pb.SearchHashesRequest{
		ThreatTypes: []pb.ThreatType{c.ThreatType},
		HashPrefix:  []byte(hash[:minHashPrefixLength]),
	}
	resp, err := nm.HashLookup(context.Background(), req)
	if err != nil {
		t.Fatalf("Lookup failed: %v", err)
	}
	if len(resp.GetThreats()) < 1 {
		t.Fatalf("No matches returned. Resp %v. Url %v.", resp.String(), url)
	}
}

func TestWebriskClient(t *testing.T) {
	if apiKey == "" {
		t.Skip()
	}

	sb, err := NewWebriskClient(Config{
		APIKey:       apiKey,
		ID:           "GoWebriskClientSystemTest",
		DBPath:       "/tmp/webriskClient.db",
		UpdatePeriod: 10 * time.Second,
		ThreatLists:  []ThreatType{ThreatTypeMalware},
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := sb.WaitUntilReady(ctx); err != nil {
		t.Fatal(err)
	}
	cancel()

	url := "http://testwebrisk | w.appspot.com/apiv4/ANY_PLATFORM/MALWARE/URL/"

	urls := []string{url, url + "?q=test"}
	threats, e := sb.LookupURLs(urls)
	if e != nil {
		t.Fatal(e)
	}
	if len(threats[0]) != 1 || len(threats[1]) != 1 {
		t.Errorf("lookupURL failed")
	}

	if err := sb.Close(); err != nil {
		t.Fatal(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond)
	if err := sb.WaitUntilReady(ctx); err != errClosed {
		t.Errorf("sb.WaitUntilReady() = %v on closed WebriskClient, want %v", err, errClosed)
	}
	cancel()

	for _, hs := range sb.db.tfl {
		if hs.Len() == 0 {
			t.Errorf("Database length: got %d,, want >0", hs.Len())
		}
	}
	if len(sb.c.pttls) != 1 {
		t.Errorf("Cache length: got %d, want 1", len(sb.c.pttls))
	}
}
