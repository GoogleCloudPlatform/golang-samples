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

func createVectorIndexes(t *testing.T, collName string,
	projectID string, dimension int32, fieldPath string) {
	dbPath := "projects/" + projectID + "/databases/(default)"
	ctx := context.Background()

	// Create client
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { client.Close() })
	collRef := client.Collection(collName)

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

func deleteIndex(t *testing.T, indexName string) {
	ctx := context.Background()

	// Create admin client
	adminClient, err := apiv1.NewFirestoreAdminClient(ctx)
	if err != nil {
		t.Logf("NewFirestoreAdminClient: %v", err)
		return
	}
	defer adminClient.Close()

	// delete index
	if err = adminClient.DeleteIndex(ctx, &adminpb.DeleteIndexRequest{
		Name: indexName,
	}); err != nil {
		t.Logf("Failed to delete index \"%s\": %+v\n", indexName, err)
	}
}

func createCoffeeBeans(t *testing.T, projectID string, collName string) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatal(err)
	}
	docs := []CoffeeBean{
		{
			Name:           "Kahawa coffee beans",
			Description:    "Information about the Kahawa coffee beans.",
			EmbeddingField: []float32{1.0, 2.0, 3.0},
			Color:          "red",
		},
		{
			Name:           "Onyx coffee beans",
			Description:    "Information about the Onyx coffee beans.",
			EmbeddingField: []float32{4.0, 5.0, 6.0},
			Color:          "brown",
		},
	}

	docRefs := []*firestore.DocumentRef{}
	for _, doc := range docs {
		ref := client.Collection(collName).NewDoc()
		docRefs = append(docRefs, ref)
		if _, err = ref.Set(ctx, doc); err != nil {
			t.Errorf("failed to upsert: %v", err)
		}
		t.Cleanup(func() {
			_, err := ref.Delete(ctx)
			if err != nil {
				t.Errorf("An error has occurred: %s", err)
			}
		})
	}
}
func TestVectorSearchBasic(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("Skipping firestore test. Set GOLANG_SAMPLES_FIRESTORE_PROJECT.")
	}

	collName := "coffee-beans"
	createCoffeeBeans(t, projectID, collName)
	createVectorIndexes(t, collName, projectID, 3, "embedding_field")

	buf := new(bytes.Buffer)
	if err := vectorSearchBasic(buf, projectID); err != nil {
		t.Errorf("vectorSearchBasic: %v", err)
	}
}
