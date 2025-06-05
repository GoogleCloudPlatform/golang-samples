// Copyright 2025 Google LLC
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
	"fmt"
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

func TestStreamingAndBatching(t *testing.T) {
	ctx := context.Background()
	project := os.Getenv("GOLANG_SAMPLES_BIGTABLE_PROJECT")
	instance := os.Getenv("GOLANG_SAMPLES_BIGTABLE_INSTANCE")
	if project == "" || instance == "" {
		t.Skip("Skipping bigtable integration test. Set GOLANG_SAMPLES_BIGTABLE_PROJECT and GOLANG_SAMPLES_BIGTABLE_INSTANCE.")
	}

	// Create client
	adminClient, err := bigtable.NewAdminClient(ctx, project, instance)
	if err != nil {
		t.Fatalf("bigtable.NewAdminClient: %v", err)
	}

	// Create table
	tableName := "mobile-time-series-" + uuid.New().String()[:8]
	if err := adminClient.DeleteTable(ctx, tableName); err != nil && status.Code(err) != codes.NotFound {
		t.Fatalf("adminClient.DeleteTable(%q): %v", tableName, err)
	}
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

	// Create column family
	columnFamilyName := "cell_plan"
	if err := adminClient.CreateColumnFamily(ctx, tableName, columnFamilyName); err != nil {
		t.Fatalf("CreateColumnFamily(%s): %v", columnFamilyName, err)
	}

	// Write test data
	client, err := bigtable.NewClient(ctx, project, instance)
	if err != nil {
		t.Fatalf("bigtable.NewClient: %v", err)
	}
	defer client.Close()
	tbl := client.Open(tableName)
	rowKeys := []string{"phone#4c410523#20190501", "phone#5c10102#20190501"}
	muts := make([]*bigtable.Mutation, len(rowKeys))
	for i := range rowKeys {
		mut := bigtable.NewMutation()
		mut.Set(columnFamilyName, "data_plan_01gb", 0, []byte("true"))
		mut.Set(columnFamilyName, "data_plan_05gb", 0, []byte("true"))
		muts[i] = mut
	}
	if _, err := tbl.ApplyBulk(ctx, rowKeys, muts); err != nil {
		t.Fatalf("tbl.ApplyBulk: %v", err)
	}

	// Run the function
	buf := new(bytes.Buffer)
	if err := streamingAndBatching(buf, project, instance, tableName); err != nil {
		t.Errorf("streamingAndBatching failed: %v", err)
	}

	// Verify the output message
	if got, want := buf.String(), "Successfully deleted cells from all rows"; !strings.Contains(got, want) {
		t.Errorf("streamingAndBatching output got %q, want substring %q", got, want)
	}

	// Verify that the cells were deleted
	for _, rowKey := range rowKeys {
		row, err := tbl.ReadRow(ctx, rowKey)
		if err != nil {
			t.Fatalf("tbl.ReadRow(%q): %v", rowKey, err)
		}
		if _, ok := row[columnFamilyName]; !ok {
			t.Errorf("row %q has no column family %q", rowKey, columnFamilyName)
			continue
		}
		for _, item := range row[columnFamilyName] {
			if item.Column == fmt.Sprintf("%s:data_plan_01gb", columnFamilyName) {
				t.Errorf("row %q still has cell %s:data_plan_01gb", rowKey, columnFamilyName)
			}
		}
	}
}
