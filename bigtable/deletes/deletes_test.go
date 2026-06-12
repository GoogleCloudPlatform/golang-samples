// Copyright 2026 Google LLC
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

package deletes

import (
	"bytes"
	"context"

	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeletes(t *testing.T) {

	ctx := context.Background()
	project := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	instance := os.Getenv("GOLANG_SAMPLES_BIGTABLE_INSTANCE")
	if project == "" || instance == "" {
		t.Skip("Skipping bigtable integration test. Set GOLANG_SAMPLES_PROJECT_ID and GOLANG_SAMPLES_BIGTABLE_INSTANCE.")
	}

	// Ensure the Bigtable instance exists.
	ensureInstance(t, ctx, project, instance)

	adminClient, err := bigtable.NewAdminClient(ctx, project, instance)
	if err != nil {
		t.Skipf("bigtable.NewAdminClient: %v", err)
	}
	defer adminClient.Close()

	tableName := "mobile-time-series-" + uuid.New().String()[:8]

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := adminClient.CreateTable(ctx, tableName); err != nil {
			if status.Code(err) == codes.AlreadyExists {
				adminClient.DeleteTable(ctx, tableName)
				time.Sleep(5 * time.Second)
			}
			r.Errorf("Could not create table %s: %v", tableName, err)
		}
	})
	if t.Failed() {
		return
	}
	defer adminClient.DeleteTable(ctx, tableName)

	if err := adminClient.CreateColumnFamily(ctx, tableName, "cell_plan"); err != nil {
		t.Fatalf("CreateColumnFamily(cell_plan): %v", err)
	}
	if err := adminClient.CreateColumnFamily(ctx, tableName, "stats_summary"); err != nil {
		t.Fatalf("CreateColumnFamily(stats_summary): %v", err)
	}

	client, err := bigtable.NewClient(ctx, project, instance)
	if err != nil {
		t.Fatalf("bigtable.NewClient: %v", err)
	}
	defer client.Close()
	tbl := client.Open(tableName)

	// Helper to reset data
	resetData := func(st *testing.T) {
		ctx := st.Context()
		// Delete all rows first
		keys := []string{
			"phone#4c410523#20190501",
			"phone#4c410523#20190502",
			"phone#4c410524#20190501",
			"phone#5c10102#20190501",
		}
		for _, k := range keys {
			mut := bigtable.NewMutation()
			mut.DeleteRow()
			if err := tbl.Apply(ctx, k, mut); err != nil {
				st.Fatalf("Failed to delete row %s: %v", k, err)
			}
		}

		// Write initial data
		muts := make([]*bigtable.Mutation, len(keys))

		// phone#4c410523#20190501
		muts[0] = bigtable.NewMutation()
		muts[0].Set("cell_plan", "data_plan_01gb", 0, []byte("true"))
		muts[0].Set("cell_plan", "data_plan_05gb", 0, []byte("true"))
		muts[0].Set("stats_summary", "connected_wifi", 0, []byte("true"))

		// phone#4c410523#20190502
		muts[1] = bigtable.NewMutation()
		muts[1].Set("cell_plan", "data_plan_01gb", 0, []byte("true"))

		// phone#4c410524#20190501
		muts[2] = bigtable.NewMutation()
		muts[2].Set("cell_plan", "data_plan_01gb", 0, []byte("true"))

		// phone#5c10102#20190501
		muts[3] = bigtable.NewMutation()
		muts[3].Set("stats_summary", "connected_wifi", 0, []byte("true"))
		muts[3].Set("stats_summary", "os_build", 0, []byte("true"))

		for i, k := range keys {
			if err := tbl.Apply(ctx, k, muts[i]); err != nil {
				st.Fatalf("Failed to setup test data for %s: %v", k, err)
			}
		}
	}

	t.Run("DeleteFromColumn", func(t *testing.T) {
		resetData(t)
		buf := new(bytes.Buffer)
		if err := deleteFromColumn(buf, project, instance, tableName); err != nil {
			t.Fatalf("deleteFromColumn: %v", err)
		}
		if got, want := buf.String(), "Successfully deleted cells from column"; !strings.Contains(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
		// Verify
		row, err := tbl.ReadRow(ctx, "phone#4c410523#20190501")
		if err != nil {
			t.Fatalf("ReadRow: %v", err)
		}
		for _, item := range row["cell_plan"] {
			if item.Column == "cell_plan:data_plan_01gb" {
				t.Error("cell_plan:data_plan_01gb should have been deleted")
			}
		}
		// check that data_plan_05gb is still there
		found := false
		for _, item := range row["cell_plan"] {
			if item.Column == "cell_plan:data_plan_05gb" {
				found = true
			}
		}
		if !found {
			t.Error("cell_plan:data_plan_05gb should still exist")
		}
	})

	t.Run("DeleteFromColumnFamily", func(t *testing.T) {
		resetData(t)
		buf := new(bytes.Buffer)
		if err := deleteFromColumnFamily(buf, project, instance, tableName); err != nil {
			t.Fatalf("deleteFromColumnFamily: %v", err)
		}
		if got, want := buf.String(), "Successfully deleted cells from family"; !strings.Contains(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
		// Verify
		row, err := tbl.ReadRow(ctx, "phone#5c10102#20190501")
		if err != nil {
			t.Fatalf("ReadRow: %v", err)
		}
		if _, ok := row["stats_summary"]; ok {
			t.Error("stats_summary family should have been deleted")
		}
	})

	t.Run("DeleteFromRow", func(t *testing.T) {
		resetData(t)
		buf := new(bytes.Buffer)
		if err := deleteFromRow(buf, project, instance, tableName); err != nil {
			t.Fatalf("deleteFromRow: %v", err)
		}
		if got, want := buf.String(), "Successfully deleted row"; !strings.Contains(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
		// Verify
		row, err := tbl.ReadRow(ctx, "phone#4c410523#20190501")
		if err != nil {
			t.Fatalf("ReadRow: %v", err)
		}
		if len(row) > 0 {
			t.Error("Row should have been deleted")
		}
	})

	t.Run("StreamingAndBatching", func(t *testing.T) {
		resetData(t)
		buf := new(bytes.Buffer)
		if err := streamingAndBatching(buf, project, instance, tableName); err != nil {
			t.Fatalf("streamingAndBatching: %v", err)
		}
		// Verify
		for _, k := range []string{"phone#4c410523#20190501", "phone#4c410523#20190502", "phone#4c410524#20190501"} {
			row, err := tbl.ReadRow(ctx, k)
			if err != nil {
				t.Fatalf("ReadRow: %v", err)
			}
			for _, item := range row["cell_plan"] {
				if item.Column == "cell_plan:data_plan_01gb" {
					t.Errorf("row %s still has cell_plan:data_plan_01gb", k)
				}
			}
		}
	})

	t.Run("DropRowRange", func(t *testing.T) {
		resetData(t)
		buf := new(bytes.Buffer)
		if err := dropRowRange(buf, project, instance, tableName); err != nil {
			t.Fatalf("dropRowRange: %v", err)
		}
		if got, want := buf.String(), "Successfully dropped row range"; !strings.Contains(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
		// Verify
		// phone#4c410523#20190501 and phone#4c410523#20190502 should be gone
		for _, k := range []string{"phone#4c410523#20190501", "phone#4c410523#20190502"} {
			row, err := tbl.ReadRow(ctx, k)
			if err != nil {
				t.Fatalf("ReadRow: %v", err)
			}
			if len(row) > 0 {
				t.Errorf("Row %s should have been dropped", k)
			}
		}
		// phone#4c410524#20190501 should still be there
		row, err := tbl.ReadRow(ctx, "phone#4c410524#20190501")
		if err != nil {
			t.Fatalf("ReadRow: %v", err)
		}
		if len(row) == 0 {
			t.Error("phone#4c410524#20190501 should NOT have been dropped")
		}
	})
}

func ensureInstance(t *testing.T, ctx context.Context, project, instance string) {
	instanceAdminClient, err := bigtable.NewInstanceAdminClient(ctx, project)
	if err != nil {
		t.Fatalf("bigtable.NewInstanceAdminClient: %v", err)
	}
	t.Cleanup(func() { instanceAdminClient.Close() })

	if _, err := instanceAdminClient.InstanceInfo(ctx, instance); err != nil {
		if status.Code(err) == codes.NotFound {
			zone := os.Getenv("GOLANG_SAMPLES_BIGTABLE_ZONE")
			if zone == "" {
				zone = "us-central1-b"
			}
			clusterID := instance + "-c1"
			instanceConf := &bigtable.InstanceConf{
				InstanceId:   instance,
				DisplayName:  instance,
				ClusterId:    clusterID,
				NumNodes:     0,
				InstanceType: bigtable.DEVELOPMENT,
				StorageType:  bigtable.SSD,
				Zone:         zone,
			}
			if err := instanceAdminClient.CreateInstance(ctx, instanceConf); err != nil {
				t.Fatalf("CreateInstance: %v", err)
			}
			t.Cleanup(func() {
				if err := instanceAdminClient.DeleteInstance(context.Background(), instance); err != nil {
					t.Errorf("DeleteInstance: %v", err)
				}
			})
		} else {
			t.Fatalf("InstanceInfo check failed: %v", err)
		}
	}
}
