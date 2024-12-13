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
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	apiv1 "cloud.google.com/go/firestore/apiv1/admin"
	"cloud.google.com/go/firestore/apiv1/admin/adminpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
)

var projectID string

func TestMain(m *testing.M) {
	ctx := context.Background()
	projectID = os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	bw := client.BulkWriter(ctx)
	colName := "users"

	docs := []struct {
		shortName string
		birthYear int
	}{
		{shortName: "aturing", birthYear: 1912},
		{shortName: "alovelace", birthYear: 1815},
		{shortName: "cbabbage", birthYear: 1791},
		{shortName: "ghopper", birthYear: 1906},
	}
	var refs []*firestore.DocumentRef

	for _, d := range docs {
		ref := client.Collection(colName).Doc(d.shortName)
		_, err := bw.Create(ref, map[string]interface{}{"birthYear": d.birthYear})
		if err != nil {
			log.Fatal(err)
		}
		refs = append(refs, ref)
	}
	bw.End()

	vsCleanup := vectorSearchSetup()
	defer vsCleanup()

	// Run the test
	m.Run()

	// New BulkWriter instance
	bw = client.BulkWriter(ctx)

	for _, d := range refs {
		_, err := bw.Delete(d)
		if err != nil {
			log.Fatal(err)
		}
	}
	bw.End()
}

func vectorSearchSetup() func() {
	vectorCollName := "coffee-beans"
	vectorQueryFieldPath := "embedding_field"
	vectorFieldDimension := int32(3)

	cleanups := []func(){}

	// Delete existing documents
	deleteTestCollection(projectID, vectorCollName)

	// Create documents
	cleanupDocs := createCoffeeBeans(projectID, vectorCollName)
	cleanups = append(cleanups, cleanupDocs)

	// Wait for single field indexes to get created
	time.Sleep(30 * time.Second)

	// Create indexes
	indexFields := [][]*adminpb.Index_IndexField{
		[]*adminpb.Index_IndexField{
			{
				FieldPath: vectorQueryFieldPath,
				ValueMode: &adminpb.Index_IndexField_VectorConfig_{
					VectorConfig: &adminpb.Index_IndexField_VectorConfig{
						Dimension: vectorFieldDimension,
						Type: &adminpb.Index_IndexField_VectorConfig_Flat{
							Flat: &adminpb.Index_IndexField_VectorConfig_FlatIndex{},
						},
					},
				},
			},
		},
		// vector indexes required for vector search with prefilter
		[]*adminpb.Index_IndexField{
			{
				FieldPath: "color",
				ValueMode: &adminpb.Index_IndexField_Order_{
					Order: adminpb.Index_IndexField_ASCENDING,
				},
			},
			{
				FieldPath: vectorQueryFieldPath,
				ValueMode: &adminpb.Index_IndexField_VectorConfig_{
					VectorConfig: &adminpb.Index_IndexField_VectorConfig{
						Dimension: vectorFieldDimension,
						Type: &adminpb.Index_IndexField_VectorConfig_Flat{
							Flat: &adminpb.Index_IndexField_VectorConfig_FlatIndex{},
						},
					},
				},
			},
		},
	}
	for _, fields := range indexFields {
		cleanup := createVectorIndex(projectID, vectorCollName, fields)
		cleanups = append(cleanups, cleanup)
	}

	return func() {
		for _, cleanup := range cleanups {
			cleanup()
		}
	}
}

func createVectorIndex(projectID, collName string, fields []*adminpb.Index_IndexField) func() {
	dbPath := "projects/" + projectID + "/databases/(default)"
	ctx := context.Background()

	// Create client
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	collRef := client.Collection(collName)

	// Create admin client
	adminClient, err := apiv1.NewFirestoreAdminClient(ctx)
	if err != nil {
		log.Fatalf("NewFirestoreAdminClient: %v", err)
	}
	defer adminClient.Close()

	indexParent := fmt.Sprintf("%s/collectionGroups/%s", dbPath, collRef.ID)

	// create vector mode indexes
	req := &adminpb.CreateIndexRequest{
		Parent: indexParent,
		Index: &adminpb.Index{
			QueryScope: adminpb.Index_COLLECTION,
			Fields:     fields,
		},
	}
	op, createErr := adminClient.CreateIndex(ctx, req)
	if createErr != nil {
		if strings.Contains(createErr.Error(), "index already exists") {
			return func() {}
		}

		log.Fatalf("CreateIndex: %v", createErr)
	}
	createdIndex, waitErr := op.Wait(ctx)
	if waitErr != nil {
		log.Fatalf("op.Wait: %v", waitErr)
	}

	return func() {
		deleteIndex(createdIndex.Name)
	}
}

func deleteIndex(indexName string) {
	ctx := context.Background()

	// Create admin client
	adminClient, err := apiv1.NewFirestoreAdminClient(ctx)
	if err != nil {
		log.Printf("NewFirestoreAdminClient: %v", err)
		return
	}
	defer adminClient.Close()

	// delete index
	if err = adminClient.DeleteIndex(ctx, &adminpb.DeleteIndexRequest{
		Name: indexName,
	}); err != nil {
		log.Printf("Failed to delete index \"%s\": %+v\n", indexName, err)
	}
}

func createCoffeeBeans(projectID string, collName string) func() {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	docs := []CoffeeBean{
		{
			Name:           "Kahawa coffee beans",
			Description:    "Information about the Kahawa coffee beans.",
			EmbeddingField: []float32{1.0, 2.0, 3.0},
			Color:          "red",
		},
		{
			Name:           "Owl coffee beans",
			Description:    "Information about the Owl coffee beans.",
			EmbeddingField: []float32{4.0, 5.0, 6.0},
			Color:          "brown",
		},
		{
			Name:           "Sleepy coffee beans",
			Description:    "Information about the Sleepy coffee beans.",
			EmbeddingField: []float32{3.0, 1.0, 2.0},
			Color:          "red",
		},
	}

	docRefs := []*firestore.DocumentRef{}
	for _, doc := range docs {
		ref := client.Collection(collName).NewDoc()
		docRefs = append(docRefs, ref)
		if _, err = ref.Set(ctx, doc); err != nil {
			log.Fatalf("failed to upsert: %v", err)
		}
	}

	return func() {
		for _, ref := range docRefs {
			testutil.RetryWithoutTest(5, 5*time.Second, func(r *testutil.R) {
				_, err := ref.Delete(ctx)
				if err != nil {
					log.Printf("Error deleting document %v: %s", ref, err)
					r.Fail()
				}
			})
		}
	}
}

func deleteTestCollection(projectID, collName string) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Delete all documents in the collName collection.
	iter := client.Collection(collName).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		testutil.RetryWithoutTest(5, 5*time.Second, func(r *testutil.R) {
			_, err = doc.Ref.Delete(ctx)
			if err != nil {
				log.Printf("Error deleting document %v: %s", doc.Ref, err)
				r.Fail()
			}
		})
	}
}
