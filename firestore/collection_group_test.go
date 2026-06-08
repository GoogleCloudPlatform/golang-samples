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

package firestore

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"cloud.google.com/go/firestore"
	apiv1 "cloud.google.com/go/firestore/apiv1/admin"
	"cloud.google.com/go/firestore/apiv1/admin/adminpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/genproto/protobuf/field_mask"
)

func TestCollectionGroup(t *testing.T) {
	tc := testutil.SystemTest(t)
	// TODO(#559): revert this to testutil.SystemTest(t).ProjectID
	// when datastore and firestore can co-exist in a project.
	projectID := getProjectID(t)

	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	collection := tc.ProjectID + "-collection-group-cities"

	// Delete all docs first to make sure collectionGroupSetup works.
	cleanupCollection(ctx, t, client, collection)

	if err := collectionGroupSetup(projectID, collection); err != nil {
		t.Fatalf("collectionGroupSetup: %v", err)
	}

	if indexCleanup, err := createIndexExemption(projectID); err != nil {
		t.Fatalf("createCollectionGroupIndex: %v", err)
	} else {
		t.Cleanup(func() {
			if indexCleanup != nil {
				err := indexCleanup()
				if err != nil {
					t.Errorf("cleanup failed: %v", err)
				}
			}
		})
	}

	buf := &bytes.Buffer{}
	if err := collectionGroupQuery(buf, projectID); err != nil {
		t.Fatalf("collectionGroupQuery: %v", err)
	}
	want := "Legion of Honor"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("collectionGroupQuery got\n----\n%s\n----\nWant to contain:\n----\n%s\n----", got, want)
	}
}

func createIndexExemption(projectID string) (func() error, error) {
	ctx := context.Background()

	// Create admin client
	adminClient, err := apiv1.NewFirestoreAdminClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewFirestoreAdminClient: %v", err)
	}
	defer adminClient.Close()

	fieldResourceName := fmt.Sprintf("projects/%s/databases/%s/collectionGroups/%s/fields/%s",
		projectID, "(default)", "landmarks", "type")

	getFieldRequest := &adminpb.GetFieldRequest{
		Name: fieldResourceName,
	}

	origField, err := adminClient.GetField(ctx, getFieldRequest)
	if err != nil {
		return nil, fmt.Errorf("GetField: %v", err)
	}

	updateReq := &adminpb.UpdateFieldRequest{
		Field: &adminpb.Field{
			Name: fieldResourceName,
			IndexConfig: &adminpb.Field_IndexConfig{
				// Providing an empty list of indexes disables all automatic indexes
				Indexes: []*adminpb.Index{
					{
						QueryScope: adminpb.Index_COLLECTION,
						Fields: []*adminpb.Index_IndexField{
							{
								FieldPath: "type",
								ValueMode: &adminpb.Index_IndexField_Order_{
									Order: adminpb.Index_IndexField_ASCENDING,
								},
							},
						},
					},
					{
						QueryScope: adminpb.Index_COLLECTION,
						Fields: []*adminpb.Index_IndexField{
							{
								FieldPath: "type",
								ValueMode: &adminpb.Index_IndexField_Order_{
									Order: adminpb.Index_IndexField_DESCENDING,
								},
							},
						},
					},
					{
						QueryScope: adminpb.Index_COLLECTION,
						Fields: []*adminpb.Index_IndexField{
							{
								FieldPath: "type",
								ValueMode: &adminpb.Index_IndexField_ArrayConfig_{
									ArrayConfig: adminpb.Index_IndexField_CONTAINS,
								},
							},
						},
					},
					{
						QueryScope: adminpb.Index_COLLECTION_GROUP,
						Fields: []*adminpb.Index_IndexField{
							{
								FieldPath: "type",
								ValueMode: &adminpb.Index_IndexField_Order_{
									Order: adminpb.Index_IndexField_ASCENDING,
								},
							},
						},
					},
				},
			},
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"index_config"},
		},
	}

	op, updateErr := adminClient.UpdateField(ctx, updateReq)
	if updateErr != nil {
		if strings.Contains(updateErr.Error(), "already exists") {
			return nil, nil
		}
		return nil, fmt.Errorf("UpdateField: %v", updateErr)
	}
	// Wait until the operation completes.
	_, waitErr := op.Wait(ctx)
	if waitErr != nil {
		return nil, fmt.Errorf("UpdateField.Wait: %v", waitErr)
	}
	return func() error {
		adminClient, err := apiv1.NewFirestoreAdminClient(ctx)
		if err != nil {
			return fmt.Errorf("NewFirestoreAdminClient: %v", err)
		}
		defer adminClient.Close()
		_, updateErr := adminClient.UpdateField(ctx, &adminpb.UpdateFieldRequest{
			Field: origField,
			UpdateMask: &field_mask.FieldMask{
				Paths: []string{"index_config"},
			},
		})
		if updateErr != nil {
			return fmt.Errorf("UpdateField: %v", updateErr)
		}
		return nil
	}, nil
}

func TestCollectionGroupPartitionQueries(t *testing.T) {
	tc := testutil.SystemTest(t)
	// TODO(#559): revert this to testutil.SystemTest(t).ProjectID
	// when datastore and firestore can co-exist in a project.
	projectID := getProjectID(t)

	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	collection := tc.ProjectID + "-collection-group-cities"

	// Delete all docs first to make sure collectionGroupSetup works.
	cleanupCollection(ctx, t, client, collection)

	if err := collectionGroupSetup(projectID, collection); err != nil {
		t.Fatalf("collectionGroupSetup: %v", err)
	}

	err = partitionQuery(ctx, client)
	if err != nil {
		t.Errorf("partitionQuery unexpected result: %s", err)
	}

	err = serializePartitionQuery(ctx, client)
	if err != nil {
		t.Errorf("serializePartitionQuery unexpected result: %s", err)
	}
}

func cleanupCollection(ctx context.Context, t *testing.T, client *firestore.Client, collection string) {
	t.Helper()
	docs, err := client.Collection(collection).DocumentRefs(ctx).GetAll()
	if err != nil {
		t.Logf("Warning: failed to get document refs for cleanup: %v", err)
		return
	}
	for _, d := range docs {
		if _, err := d.Delete(ctx); err != nil {
			t.Logf("Warning: failed to delete doc %s: %v", d.ID, err)
		}
	}
}
