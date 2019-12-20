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

package reads

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"io"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/bigtable"
)

func TestReads(t *testing.T) {
	ctx := context.Background()
	project := os.Getenv("GOLANG_SAMPLES_BIGTABLE_PROJECT")
	instance := os.Getenv("GOLANG_SAMPLES_BIGTABLE_INSTANCE")
	if project == "" || instance == "" {
		t.Skip("Skipping bigtable integration test. Set GOLANG_SAMPLES_BIGTABLE_PROJECT and GOLANG_SAMPLES_BIGTABLE_INSTANCE.")
	}
	adminClient, err := bigtable.NewAdminClient(ctx, project, instance)

	uuid, err := uuid.NewRandom()
	tableName := fmt.Sprintf("mobile-time-series-%s", uuid.String()[:8])
	adminClient.DeleteTable(ctx, tableName)

	if err := adminClient.CreateTable(ctx, tableName); err != nil {
		t.Fatalf("Could not create table %s: %v", tableName, err)
	}

	if err := adminClient.CreateColumnFamily(ctx, tableName, "stats_summary"); err != nil {
		adminClient.DeleteTable(ctx, tableName)
		t.Fatalf("CreateColumnFamily(%s): %v", "stats_summary", err)
	}

	timestamp := bigtable.Now().TruncateToMilliseconds()
	writeTestData(err, ctx, project, instance, tableName, timestamp, t)

	tests := []struct {
		name   string
		filter func(io.Writer, string, string, string) error
		want   string
	}{
		{
			name: "readRow", filter: readRow, want: fmt.Sprintf(
				`Reading data for phone#4c410523#20190501:
Column Family stats_summary
	connected_cell: 1 @%[1]d
	connected_wifi: 1 @%[1]d
	os_build: PQ2A.190405.003 @%[1]d`, timestamp),
		},
		{
			name: "readRowPartial", filter: readRowPartial, want: fmt.Sprintf(
				`Reading data for phone#4c410523#20190501:
Column Family stats_summary
	os_build: PQ2A.190405.003 @%[1]d`, timestamp),
		},
		{
			name: "readRows", filter: readRows, want: fmt.Sprintf(
				`Reading data for phone#4c410523#20190501:
Column Family stats_summary
	connected_cell: 1 @%[1]d
	connected_wifi: 1 @%[1]d
	os_build: PQ2A.190405.003 @%[1]d

Reading data for phone#4c410523#20190502:
Column Family stats_summary
	connected_cell: 1 @%[1]d
	connected_wifi: 1 @%[1]d
	os_build: PQ2A.190405.004 @%[1]d`, timestamp),
		},
		{
			name: "readRowRange", filter: readRowRange, want: fmt.Sprintf(
				`Reading data for phone#4c410523#20190501:
Column Family stats_summary
	connected_cell: 1 @%[1]d
	connected_wifi: 1 @%[1]d
	os_build: PQ2A.190405.003 @%[1]d

Reading data for phone#4c410523#20190502:
Column Family stats_summary
	connected_cell: 1 @%[1]d
	connected_wifi: 1 @%[1]d
	os_build: PQ2A.190405.004 @%[1]d

Reading data for phone#4c410523#20190505:
Column Family stats_summary
	connected_cell: 0 @%[1]d
	connected_wifi: 1 @%[1]d
	os_build: PQ2A.190406.000 @%[1]d`, timestamp),
		},
		{
			name: "readRowRanges", filter: readRowRanges, want: fmt.Sprintf(
				`Reading data for phone#4c410523#20190501:
Column Family stats_summary
	connected_cell: 1 @%[1]d
	connected_wifi: 1 @%[1]d
	os_build: PQ2A.190405.003 @%[1]d

Reading data for phone#4c410523#20190502:
Column Family stats_summary
	connected_cell: 1 @%[1]d
	connected_wifi: 1 @%[1]d
	os_build: PQ2A.190405.004 @%[1]d

Reading data for phone#4c410523#20190505:
Column Family stats_summary
	connected_cell: 0 @%[1]d
	connected_wifi: 1 @%[1]d
	os_build: PQ2A.190406.000 @%[1]d

Reading data for phone#5c10102#20190501:
Column Family stats_summary
	connected_cell: 1 @%[1]d
	connected_wifi: 1 @%[1]d
	os_build: PQ2A.190401.002 @%[1]d

Reading data for phone#5c10102#20190502:
Column Family stats_summary
	connected_cell: 1 @%[1]d
	connected_wifi: 0 @%[1]d
	os_build: PQ2A.190406.000 @%[1]d`, timestamp),
		},
		{
			name: "readPrefix", filter: readPrefix, want: fmt.Sprintf(
				`Reading data for phone#4c410523#20190501:
Column Family stats_summary
	connected_cell: 1 @%[1]d
	connected_wifi: 1 @%[1]d
	os_build: PQ2A.190405.003 @%[1]d

Reading data for phone#4c410523#20190502:
Column Family stats_summary
	connected_cell: 1 @%[1]d
	connected_wifi: 1 @%[1]d
	os_build: PQ2A.190405.004 @%[1]d

Reading data for phone#4c410523#20190505:
Column Family stats_summary
	connected_cell: 0 @%[1]d
	connected_wifi: 1 @%[1]d
	os_build: PQ2A.190406.000 @%[1]d

Reading data for phone#5c10102#20190501:
Column Family stats_summary
	connected_cell: 1 @%[1]d
	connected_wifi: 1 @%[1]d
	os_build: PQ2A.190401.002 @%[1]d

Reading data for phone#5c10102#20190502:
Column Family stats_summary
	connected_cell: 1 @%[1]d
	connected_wifi: 0 @%[1]d
	os_build: PQ2A.190406.000 @%[1]d`, timestamp),
		},
		{
			name: "readFilter", filter: readFilter, want: fmt.Sprintf(
				`Reading data for phone#4c410523#20190501:
Column Family stats_summary
	os_build: PQ2A.190405.003 @%[1]d

Reading data for phone#4c410523#20190502:
Column Family stats_summary
	os_build: PQ2A.190405.004 @%[1]d

Reading data for phone#4c410523#20190505:
Column Family stats_summary
	os_build: PQ2A.190406.000 @%[1]d

Reading data for phone#5c10102#20190501:
Column Family stats_summary
	os_build: PQ2A.190401.002 @%[1]d

Reading data for phone#5c10102#20190502:
Column Family stats_summary
	os_build: PQ2A.190406.000 @%[1]d`, timestamp),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			if err = tt.filter(buf, project, instance, tableName); err != nil {
				t.Errorf("Testing %s: %v", tt.name, err)
			}

			got := buf.String()

			if diff := cmp.Diff(tt.want, strings.TrimSpace(got)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}

	adminClient.DeleteTable(ctx, tableName)
}

func writeTestData(err error, ctx context.Context, project string, instance string, tableName string, timestamp bigtable.Timestamp, t *testing.T) {

	client, err := bigtable.NewClient(ctx, project, instance)
	tbl := client.Open(tableName)

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
