// Copyright 2020 Google LLC
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

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestPartion(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("Skipping firestore test. Set GOLANG_SAMPLES_FIRESTORE_PROJECT.")
	}
	collectionID := tc.ProjectID + "-collection"

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	docs, err := client.Collection(collectionID).DocumentRefs(ctx).GetAll()
	if err == nil {
		for _, d := range docs {
			d.Delete(ctx)
		}
	}
	// Minimum partition size is 128.
	documentCount := 128*2 + 127
	collection := client.Collection(collectionID)
	collectionGroupID := collectionID + "-subcollection"
	for i := 0; i < documentCount; i++ {
		doc := fmt.Sprintf("doc-" + strconv.Itoa(i+1))
		if _, err := collection.Doc(doc).Collection(collectionGroupID).NewDoc().Set(ctx, map[string]int{
			"id": i,
		}); err != nil {
			t.Errorf("Set: %v", err)
		}
	}
	parent := fmt.Sprintf("projects/%v/databases/(default)/documents", projectID)
	if err := partitionQuery(ioutil.Discard, parent, collectionGroupID); err != nil {
		t.Fatalf("partitionQuery: %v", err)
	}
}
