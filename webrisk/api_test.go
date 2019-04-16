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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	pb "github.com/GoogleCloudPlatform/golang-samples/webrisk/internal/webrisk_proto"
	"github.com/golang/protobuf/proto"
)

type mockAPI struct {
	listUpdate func(context.Context, *pb.ComputeThreatListDiffRequest) (*pb.ComputeThreatListDiffResponse, error)
	hashLookup func(context.Context, *pb.SearchHashesRequest) (*pb.SearchHashesResponse, error)
}

func (m *mockAPI) ListUpdate(ctx context.Context, req *pb.ComputeThreatListDiffRequest) (*pb.ComputeThreatListDiffResponse, error) {
	return m.listUpdate(ctx, req)
}

func (m *mockAPI) HashLookup(ctx context.Context, req *pb.SearchHashesRequest) (*pb.SearchHashesResponse, error) {
	return m.hashLookup(ctx, req)
}

func TestNetAPI(t *testing.T) {
	var gotReq, wantReq proto.Message
	var gotResp, wantResp proto.Message
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p []byte
		var err error
		if p, err = ioutil.ReadAll(r.Body); err != nil {
			t.Fatalf("unexpected ioutil.ReadAll error: %v", err)
		}
		if err := proto.Unmarshal(p, gotReq); err != nil {
			t.Fatalf("unexpected proto.Unmarshal error: %v", err)
		}
		if p, err = proto.Marshal(wantResp); err != nil {
			t.Fatalf("unexpected proto.Marshal error: %v", err)
		}
		if _, err := w.Write(p); err != nil {
			t.Fatalf("unexpected ResponseWriter.Write error: %v", err)
		}
	}))
	defer ts.Close()

	api, err := newNetAPI(ts.URL, "fizzbuzz", "")
	if err != nil {
		t.Errorf("unexpected newNetAPI error: %v", err)
	}

	// Test that ListUpdate marshal/unmarshal works.
	wantReq = &pb.ComputeThreatListDiffRequest{
		ThreatType: pb.ThreatType_MALWARE,
		Constraints: &pb.ComputeThreatListDiffRequest_Constraints{
			SupportedCompressions: []pb.CompressionType{1, 2, 3}},
	}
	wantResp = &pb.ComputeThreatListDiffResponse{
		ResponseType: 1,
		Checksum:     &pb.ComputeThreatListDiffResponse_Checksum{Sha256: []byte("abcd")},
		Removals: &pb.ThreatEntryRemovals{
			RawIndices: &pb.RawIndices{Indices: []int32{1, 2, 3}},
		},
	}
	gotReq = &pb.ComputeThreatListDiffRequest{}
	resp1, err := api.ListUpdate(context.Background(), wantReq.(*pb.ComputeThreatListDiffRequest))
	gotResp = resp1
	if err != nil {
		t.Errorf("unexpected ListUpdate error: %v", err)
	}
	if !reflect.DeepEqual(gotReq, wantReq) {
		t.Errorf("mismatching ListUpdate requests:\ngot  %+v\nwant %+v", gotReq, wantReq)
	}
	if !reflect.DeepEqual(gotResp, wantResp) {
		t.Errorf("mismatching ListUpdate responses:\ngot  %+v\nwant %+v", gotResp, wantResp)
	}

	// Test that HashLookup marshal/unmarshal works.
	wantReq = &pb.SearchHashesRequest{
		HashPrefix:  []byte("aaaa"),
		ThreatTypes: []pb.ThreatType{1, 2, 3},
	}
	wantResp = &pb.SearchHashesResponse{Threats: []*pb.SearchHashesResponse_ThreatHash{{
		ThreatTypes: []pb.ThreatType{pb.ThreatType_MALWARE},
		Hash:        []byte("abcd")}}}
	gotReq = &pb.SearchHashesRequest{}
	resp2, err := api.HashLookup(context.Background(), wantReq.(*pb.SearchHashesRequest))
	gotResp = resp2
	if err != nil {
		t.Errorf("unexpected HashLookup error: %v", err)
	}
	if !reflect.DeepEqual(gotReq, wantReq) {
		t.Errorf("mismatching HashLookup requests:\ngot  %+v\nwant %+v", gotReq, wantReq)
	}
	if !reflect.DeepEqual(gotResp, wantResp) {
		t.Errorf("mismatching HashLookup responses:\ngot  %+v\nwant %+v", gotResp, wantResp)
	}

	// Test canceled Context returns an error.
	wantReq = &pb.SearchHashesRequest{
		HashPrefix:  []byte("aaaa"),
		ThreatTypes: []pb.ThreatType{1, 2, 3},
	}
	wantResp = &pb.SearchHashesResponse{Threats: []*pb.SearchHashesResponse_ThreatHash{{
		ThreatTypes: []pb.ThreatType{pb.ThreatType_MALWARE},
		Hash:        []byte("abcd")},
	}}
	gotReq = &pb.SearchHashesRequest{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = api.HashLookup(ctx, wantReq.(*pb.SearchHashesRequest))
	if err == nil {
		t.Errorf("unexpected HashLookup success, wanted HTTP request canceled")
	}
}
