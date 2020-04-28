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
	"github.com/google/uuid"
)

func TestBigtableRead(t *testing.T) {
	ctx := context.Background()
	project := os.Getenv("GOLANG_SAMPLES_BIGTABLE_PROJECT")
	instance := os.Getenv("GOLANG_SAMPLES_BIGTABLE_INSTANCE")
	if project == "" || instance == "" {
		t.Skip("Skipping functions bigtable integration test. Set GOLANG_SAMPLES_BIGTABLE_PROJECT and GOLANG_SAMPLES_BIGTABLE_INSTANCE.")
	}

	adminClient, err := bigtable.NewAdminClient(ctx, project, instance)

	uuid, err := uuid.NewRandom()
	tableId := fmt.Sprintf("mobile-time-series-%s", uuid.String()[:8])
	adminClient.DeleteTable(ctx, tableId)

	if err := adminClient.CreateTable(ctx, tableId); err != nil {
		t.Fatalf("Could not create table %s: %v", tableId, err)
	}

	if err := adminClient.CreateColumnFamily(ctx, tableId, "stats_summary"); err != nil {
		adminClient.DeleteTable(ctx, tableId)
		t.Fatalf("CreateColumnFamily(%s): %v", "stats_summary", err)
	}

	timestamp := bigtable.Now().TruncateToMilliseconds()
	writeTestData(err, ctx, project, instance, tableId, timestamp, t)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("projectId", project)
	req.Header.Set("instanceId", instance)
	req.Header.Set("tableId", tableId)
	rr := httptest.NewRecorder()
	BigtableRead(rr, req)

	t.Log(rr.Body.String())
	got := (rr.Body.String())
	if want :=
		`Rowkey: phone#4c410523#20190501, os_build:  PQ2A.190405.003
Rowkey: phone#4c410523#20190502, os_build:  PQ2A.190405.004
Rowkey: phone#4c410523#20190505, os_build:  PQ2A.190406.000
Rowkey: phone#5c10102#20190501, os_build:  PQ2A.190401.002
Rowkey: phone#5c10102#20190502, os_build:  PQ2A.190406.000`; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}

	if rr.Code != http.StatusOK {
		t.Fatalf("BigtableRead got code %v, want %v", rr.Code, http.StatusOK)
	}

	adminClient.DeleteTable(ctx, tableId)
}

func writeTestData(err error, ctx context.Context, project string, instance string, tableId string, timestamp bigtable.Timestamp, t *testing.T) {

	client, err := bigtable.NewClient(ctx, project, instance)
	tbl := client.Open(tableId)

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
