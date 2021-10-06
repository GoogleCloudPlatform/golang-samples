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

package writes

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestWrites(t *testing.T) {
	tc := testutil.SystemTest(t)

	ctx := context.Background()
	project := os.Getenv("GOLANG_SAMPLES_BIGTABLE_PROJECT")
	instance := os.Getenv("GOLANG_SAMPLES_BIGTABLE_INSTANCE")
	if project == "" || instance == "" {
		t.Skip("Skipping bigtable integration test. Set GOLANG_SAMPLES_BIGTABLE_PROJECT and GOLANG_SAMPLES_BIGTABLE_INSTANCE.")
	}
	adminClient, err := bigtable.NewAdminClient(ctx, project, instance)
	if err != nil {
		t.Skipf("bigtable.NewAdminClient: %v", err)
	}

	tableName := "mobile-time-series-" + tc.ProjectID
	adminClient.DeleteTable(ctx, tableName)

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := adminClient.CreateTable(ctx, tableName); err != nil {
			// Just in case the table exists, try to delete it again.
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

	columnFamilyName := "stats_summary"
	if err := adminClient.CreateColumnFamily(ctx, tableName, columnFamilyName); err != nil {
		adminClient.DeleteTable(ctx, tableName)
		t.Fatalf("CreateColumnFamily(%s): %v", columnFamilyName, err)
	}

	buf := new(bytes.Buffer)
	if err = writeSimple(buf, project, instance, tableName); err != nil {
		t.Errorf("TestWriteSimple: %v", err)
	}

	if got, want := buf.String(), "Successfully wrote row"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}

	buf.Reset()
	if err = writeConditionally(buf, project, instance, tableName); err != nil {
		t.Errorf("TestWriteConditionally: %v", err)
	}

	if got, want := buf.String(), "Successfully updated row's os_name"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}

	buf.Reset()
	if err = writeIncrement(buf, project, instance, tableName); err != nil {
		t.Errorf("TestWriteIncrement: %v", err)
	}

	if got, want := buf.String(), "Successfully updated row"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}

	buf.Reset()
	if err = writeBatch(buf, project, instance, tableName); err != nil {
		t.Errorf("TestWriteBatch: %v", err)
	}

	if got, want := buf.String(), "Successfully wrote 2 rows"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}

	adminClient.DeleteTable(ctx, tableName)
}
