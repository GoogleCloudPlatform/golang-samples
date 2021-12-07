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
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCollectionGroup(t *testing.T) {
	tc := testutil.SystemTest(t)
	// TODO(#559): revert this to testutil.SystemTest(t).ProjectID
	// when datastore and firestore can co-exist in a project.
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("Skipping firestore test. Set GOLANG_SAMPLES_FIRESTORE_PROJECT.")
	}

	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	collection := tc.ProjectID + "-collection-group-cities"

	// Delete all docs first to make sure collectionGroupSetup works.
	docs, err := client.Collection(collection).DocumentRefs(ctx).GetAll()
	// Ignore errors; this isn't essential.
	if err == nil {
		for _, d := range docs {
			d.Delete(ctx)
		}
	}

	if err := collectionGroupSetup(projectID, collection); err != nil {
		t.Fatalf("collectionGroupSetup: %v", err)
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

func TestCollectionGroupPartitionQueries(t *testing.T) {
	tc := testutil.SystemTest(t)
	// TODO(#559): revert this to testutil.SystemTest(t).ProjectID
	// when datastore and firestore can co-exist in a project.
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("Skipping firestore test. Set GOLANG_SAMPLES_FIRESTORE_PROJECT.")
	}

	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	collection := tc.ProjectID + "-collection-group-cities"

	// Delete all docs first to make sure collectionGroupSetup works.
	docs, err := client.Collection(collection).DocumentRefs(ctx).GetAll()
	// Ignore errors; this isn't essential.
	if err == nil {
		for _, d := range docs {
			d.Delete(ctx)
		}
	}

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
