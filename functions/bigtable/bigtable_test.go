// Copyright 2020 Google LLC
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

package bigtable

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/bigtable"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

func TestBigtableRead(t *testing.T) {
	ctx := context.Background()
	testutil.SystemTest(t)
	projectID := os.Getenv("GOLANG_SAMPLES_BIGTABLE_PROJECT")
	instanceID := os.Getenv("GOLANG_SAMPLES_BIGTABLE_INSTANCE")
	if projectID == "" || instanceID == "" {
		t.Skip("Skipping functions bigtable integration test. Set GOLANG_SAMPLES_BIGTABLE_PROJECT and GOLANG_SAMPLES_BIGTABLE_INSTANCE.")
	}

	adminClient, err := bigtable.NewAdminClient(ctx, projectID, instanceID)
	if err != nil {
		t.Fatalf("bigtable.NewAdminClient: %v", err)
	}

	uuid, _ := uuid.NewRandom()
	tableID := fmt.Sprintf("mobile-time-series-%s", uuid.String()[:8])
	adminClient.DeleteTable(ctx, tableID)

	if err := adminClient.CreateTable(ctx, tableID); err != nil {
		t.Fatalf("adminClient.CreateTable %s: %v", tableID, err)
	}

	if err := adminClient.CreateColumnFamily(ctx, tableID, "stats_summary"); err != nil {
		adminClient.DeleteTable(ctx, tableID)
		t.Fatalf("adminClient.CreateColumnFamily(%s): %v", "stats_summary", err)
	}

	timestamp := bigtable.Now().TruncateToMilliseconds()
	writeTestData(ctx, projectID, instanceID, tableID, timestamp, t)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("projectId", projectID)
	req.Header.Set("instanceId", instanceID)
	req.Header.Set("tableID", tableID)
	rr := httptest.NewRecorder()
	BigtableRead(rr, req)

	want :=
		`Rowkey: phone#4c410523#20190501, os_build:  PQ2A.190405.003
Rowkey: phone#4c410523#20190502, os_build:  PQ2A.190405.004
Rowkey: phone#4c410523#20190505, os_build:  PQ2A.190406.000
Rowkey: phone#5c10102#20190501, os_build:  PQ2A.190401.002
Rowkey: phone#5c10102#20190502, os_build:  PQ2A.190406.000`
	if got := rr.Body.String(); !strings.Contains(got, want) {
		t.Errorf("TestBigtableRead(): got %q, want %q", got, want)
	}

	if rr.Code != http.StatusOK {
		t.Fatalf("BigtableRead got code %v, want %v", rr.Code, http.StatusOK)
	}

	adminClient.DeleteTable(ctx, tableID)
}

func writeTestData(ctx context.Context, project string, instance string, tableID string, timestamp bigtable.Timestamp, t *testing.T) {

	client, _ := bigtable.NewClient(ctx, project, instance)
	tbl := client.Open(tableID)

	var muts []*bigtable.Mutation
	rowKeys := []string{
		"phone#4c410523#20190501",
		"phone#4c410523#20190502",
		"phone#4c410523#20190505",
		"phone#5c10102#20190501",
		"phone#5c10102#20190502",
	}

	mut := bigtable.NewMutation()
	mut.Set("stats_summary", "connected_cell", timestamp, []byte("1"))
	mut.Set("stats_summary", "connected_wifi", timestamp, []byte("1"))
	mut.Set("stats_summary", "os_build", timestamp, []byte("PQ2A.190405.003"))
	muts = append(muts, mut)
	mut = bigtable.NewMutation()
	mut.Set("stats_summary", "connected_cell", timestamp, []byte("1"))
	mut.Set("stats_summary", "connected_wifi", timestamp, []byte("1"))
	mut.Set("stats_summary", "os_build", timestamp, []byte("PQ2A.190405.004"))
	muts = append(muts, mut)
	mut = bigtable.NewMutation()
	mut.Set("stats_summary", "connected_cell", timestamp, []byte("0"))
	mut.Set("stats_summary", "connected_wifi", timestamp, []byte("1"))
	mut.Set("stats_summary", "os_build", timestamp, []byte("PQ2A.190406.000"))
	muts = append(muts, mut)
	mut = bigtable.NewMutation()
	mut.Set("stats_summary", "connected_cell", timestamp, []byte("1"))
	mut.Set("stats_summary", "connected_wifi", timestamp, []byte("1"))
	mut.Set("stats_summary", "os_build", timestamp, []byte("PQ2A.190401.002"))
	muts = append(muts, mut)
	mut = bigtable.NewMutation()
	mut.Set("stats_summary", "connected_cell", timestamp, []byte("1"))
	mut.Set("stats_summary", "connected_wifi", timestamp, []byte("0"))
	mut.Set("stats_summary", "os_build", timestamp, []byte("PQ2A.190406.000"))
	muts = append(muts, mut)

	if _, err := tbl.ApplyBulk(ctx, rowKeys, muts); err != nil {
		t.Errorf("ApplyBulk: %v", err)
	}
}
