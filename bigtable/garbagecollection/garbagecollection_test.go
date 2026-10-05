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

package garbagecollection

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
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
)

func assertGCRule(ctx context.Context, t *testing.T, adminClient *admin.BigtableTableAdminClient, tablePath, columnFamily string, want *adminpb.GcRule) {
	t.Helper()
	tbl, err := adminClient.GetTable(ctx, &adminpb.GetTableRequest{
		Name: tablePath,
		View: adminpb.Table_SCHEMA_VIEW,
	})
	if err != nil {
		t.Fatalf("GetTable(%s): %v", tablePath, err)
	}
	cf, ok := tbl.GetColumnFamilies()[columnFamily]
	if !ok {
		t.Fatalf("column family %q not found in table %s", columnFamily, tablePath)
	}
	if got := cf.GetGcRule(); !proto.Equal(got, want) {
		t.Errorf("column family %q GC rule: got %v, want %v", columnFamily, got, want)
	}
}

func TestGarbageCollection(t *testing.T) {
	tc := testutil.SystemTest(t)

	ctx := context.Background()
	project := os.Getenv("GOLANG_SAMPLES_BIGTABLE_PROJECT")
	instance := os.Getenv("GOLANG_SAMPLES_BIGTABLE_INSTANCE")
	if project == "" || instance == "" {
		t.Skip("Skipping bigtable integration test. Set GOLANG_SAMPLES_BIGTABLE_PROJECT and GOLANG_SAMPLES_BIGTABLE_INSTANCE.")
	}
	adminClient, err := admin.NewBigtableTableAdminClient(ctx)
	if err != nil {
		t.Fatalf("admin.NewBigtableTableAdminClient: %v", err)
	}
	defer adminClient.Close()

	uuid, err := uuid.NewRandom()
	tableName := fmt.Sprintf("gc-table-%s-%s", tc.ProjectID, uuid.String()[:8])
	instancePath := fmt.Sprintf("projects/%s/instances/%s", project, instance)
	tablePath := fmt.Sprintf("%s/tables/%s", instancePath, tableName)
	adminClient.DeleteTable(ctx, &adminpb.DeleteTableRequest{Name: tablePath})

	createTableReq := &adminpb.CreateTableRequest{
		Parent:  instancePath,
		TableId: tableName,
		Table:   &adminpb.Table{},
	}
	if _, err := adminClient.CreateTable(ctx, createTableReq); err != nil {
		t.Fatalf("Could not create table %s: %v", tableName, err)
	}
	defer adminClient.DeleteTable(ctx, &adminpb.DeleteTableRequest{Name: tablePath})

	maxAge := durationpb.New(time.Hour * 24 * 5)

	buf := new(bytes.Buffer)
	if err = createFamilyGCMaxAge(buf, project, instance, tableName); err != nil {
		t.Errorf("TestGarbageCollection: %v", err)
	}

	got := buf.String()
	if want := "created column family cf1 with policy:"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
	assertGCRule(ctx, t, adminClient, tablePath, "cf1", &adminpb.GcRule{
		Rule: &adminpb.GcRule_MaxAge{
			MaxAge: maxAge,
		},
	})

	buf.Reset()
	if err = updateGCRule(buf, project, instance, tableName); err != nil {
		t.Errorf("TestGarbageCollection: %v", err)
	}

	got = buf.String()
	if want := "Updated column family cf1 GC rule with policy:"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
	assertGCRule(ctx, t, adminClient, tablePath, "cf1", &adminpb.GcRule{
		Rule: &adminpb.GcRule_MaxNumVersions{
			MaxNumVersions: 1,
		},
	})

	buf.Reset()
	if err = createFamilyGCMaxVersions(buf, project, instance, tableName); err != nil {
		t.Errorf("TestGarbageCollection: %v", err)
	}

	got = buf.String()
	if want := "created column family cf2 with policy:"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
	assertGCRule(ctx, t, adminClient, tablePath, "cf2", &adminpb.GcRule{
		Rule: &adminpb.GcRule_MaxNumVersions{
			MaxNumVersions: 2,
		},
	})

	buf.Reset()
	if err = createFamilyGCUnion(buf, project, instance, tableName); err != nil {
		t.Errorf("TestGarbageCollection: %v", err)
	}

	got = buf.String()
	if want := "created column family cf3 with policy:"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
	assertGCRule(ctx, t, adminClient, tablePath, "cf3", &adminpb.GcRule{
		Rule: &adminpb.GcRule_Union_{
			Union: &adminpb.GcRule_Union{
				Rules: []*adminpb.GcRule{
					{
						Rule: &adminpb.GcRule_MaxNumVersions{
							MaxNumVersions: 2,
						},
					},
					{
						Rule: &adminpb.GcRule_MaxAge{
							MaxAge: maxAge,
						},
					},
				},
			},
		},
	})

	buf.Reset()
	if err = createFamilyGCIntersect(buf, project, instance, tableName); err != nil {
		t.Errorf("TestGarbageCollection: %v", err)
	}

	got = buf.String()
	if want := "created column family cf4 with policy:"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
	assertGCRule(ctx, t, adminClient, tablePath, "cf4", &adminpb.GcRule{
		Rule: &adminpb.GcRule_Intersection_{
			Intersection: &adminpb.GcRule_Intersection{
				Rules: []*adminpb.GcRule{
					{
						Rule: &adminpb.GcRule_MaxNumVersions{
							MaxNumVersions: 2,
						},
					},
					{
						Rule: &adminpb.GcRule_MaxAge{
							MaxAge: maxAge,
						},
					},
				},
			},
		},
	})

	buf.Reset()
	if err = createFamilyGCNested(buf, project, instance, tableName); err != nil {
		t.Errorf("TestGarbageCollection: %v", err)
	}

	got = buf.String()
	if want := "created column family cf5 with policy:"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
	assertGCRule(ctx, t, adminClient, tablePath, "cf5", &adminpb.GcRule{
		Rule: &adminpb.GcRule_Union_{
			Union: &adminpb.GcRule_Union{
				Rules: []*adminpb.GcRule{
					{
						Rule: &adminpb.GcRule_MaxNumVersions{
							MaxNumVersions: 10,
						},
					},
					{
						Rule: &adminpb.GcRule_Intersection_{
							Intersection: &adminpb.GcRule_Intersection{
								Rules: []*adminpb.GcRule{
									{
										Rule: &adminpb.GcRule_MaxNumVersions{
											MaxNumVersions: 2,
										},
									},
									{
										Rule: &adminpb.GcRule_MaxAge{
											MaxAge: maxAge,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})
}
