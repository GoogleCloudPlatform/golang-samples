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
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	admin "cloud.google.com/go/bigtable/admin/apiv2"
	"cloud.google.com/go/bigtable/admin/apiv2/adminpb"
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
	adminClient, err := admin.NewBigtableTableAdminClient(ctx)
	if err != nil {
		t.Skipf("admin.NewBigtableTableAdminClient: %v", err)
	}
	defer adminClient.Close()

	tableName := "mobile-time-series-" + tc.ProjectID
	instancePath := fmt.Sprintf("projects/%s/instances/%s", project, instance)
	tablePath := fmt.Sprintf("%s/tables/%s", instancePath, tableName)
	adminClient.DeleteTable(ctx, &adminpb.DeleteTableRequest{Name: tablePath})

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if _, err := adminClient.CreateTable(ctx, &adminpb.CreateTableRequest{
			Parent:  instancePath,
			TableId: tableName,
			Table: &adminpb.Table{
				ColumnFamilies: map[string]*adminpb.ColumnFamily{
					"stats_summary": {},
					"view_count": {
						ValueType: &adminpb.Type{
							Kind: &adminpb.Type_AggregateType{
								AggregateType: &adminpb.Type_Aggregate{
									InputType: &adminpb.Type{
										Kind: &adminpb.Type_Int64Type{
											Int64Type: &adminpb.Type_Int64{
												Encoding: &adminpb.Type_Int64_Encoding{
													Encoding: &adminpb.Type_Int64_Encoding_BigEndianBytes_{
														BigEndianBytes: &adminpb.Type_Int64_Encoding_BigEndianBytes{},
													},
												},
											},
										},
									},
									Aggregator: &adminpb.Type_Aggregate_Sum_{
										Sum: &adminpb.Type_Aggregate_Sum{},
									},
								},
							},
						},
					},
				},
			},
		}); err != nil {
			// Just in case the table exists, try to delete it again.
			if status.Code(err) == codes.AlreadyExists {
				adminClient.DeleteTable(ctx, &adminpb.DeleteTableRequest{Name: tablePath})
				time.Sleep(5 * time.Second)
			}
			r.Errorf("Could not create table %s: %v", tableName, err)
		}
	})
	if t.Failed() {
		return
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
	if err = writeAggregate(buf, project, instance, tableName); err != nil {
		t.Errorf("TestWriteAggregate: %v", err)
	}

	buf.Reset()
	if err = writeBatch(buf, project, instance, tableName); err != nil {
		t.Errorf("TestWriteBatch: %v", err)
	}

	if got, want := buf.String(), "Successfully wrote 2 rows"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}

	adminClient.DeleteTable(ctx, &adminpb.DeleteTableRequest{Name: tablePath})
}
