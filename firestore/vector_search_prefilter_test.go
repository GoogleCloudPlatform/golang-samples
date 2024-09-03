// Copyright 2024 Google LLC
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
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/firestore"
	apiv1 "cloud.google.com/go/firestore/apiv1/admin"
	"cloud.google.com/go/firestore/apiv1/admin/adminpb"
)

func createPrefilterVectorIndexes(t *testing.T, collName string,
	projectID string, dimension int32, fieldPath string) {
	dbPath := "projects/" + projectID + "/databases/(default)"
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatal(err)
	}
	collRef := client.Collection(collName)
	t.Cleanup(func() { client.Close() })

	// Create admin client
	adminClient, err := apiv1.NewFirestoreAdminClient(ctx)
	if err != nil {
		t.Fatalf("NewFirestoreAdminClient: %v", err)
	}
	t.Cleanup(func() { adminClient.Close() })

	indexParent := fmt.Sprintf("%s/collectionGroups/%s", dbPath, collRef.ID)

	// create vector mode indexes
	req := &adminpb.CreateIndexRequest{
		Parent: indexParent,
		Index: &adminpb.Index{
			QueryScope: adminpb.Index_COLLECTION,
			Fields: []*adminpb.Index_IndexField{
				{
					FieldPath: "color",
					ValueMode: &adminpb.Index_IndexField_Order_{
						Order: adminpb.Index_IndexField_ASCENDING,
					},
				},
				{
					FieldPath: fieldPath,
					ValueMode: &adminpb.Index_IndexField_VectorConfig_{
						VectorConfig: &adminpb.Index_IndexField_VectorConfig{
							Dimension: dimension,
							Type: &adminpb.Index_IndexField_VectorConfig_Flat{
								Flat: &adminpb.Index_IndexField_VectorConfig_FlatIndex{},
							},
						},
					},
				},
			},
		},
	}
	op, createErr := adminClient.CreateIndex(ctx, req)
	if createErr != nil {
		if strings.Contains(createErr.Error(), "index already exists") {
			return
		}
		t.Fatalf("CreateIndex: %v", createErr)
	}

	createdIndex, waitErr := op.Wait(ctx)
	if waitErr != nil {
		t.Fatalf("op.Wait: %v", waitErr)
	}
	t.Cleanup(func() { deleteIndex(t, createdIndex.Name) })
}
func TestVectorSearchPrefilter(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("Skipping firestore test. Set GOLANG_SAMPLES_FIRESTORE_PROJECT.")
	}

	collName := "coffee-beans"
	createCoffeeBeans(t, projectID, collName)
	createPrefilterVectorIndexes(t, collName, projectID, 3, "embedding_field")

	buf := new(bytes.Buffer)
	if err := vectorSearchPrefilter(buf, projectID); err != nil {
		t.Errorf("vectorSearchBasic: %v", err)
	}
}
