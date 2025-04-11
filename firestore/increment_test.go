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
	"context"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestUpdateDocumentIncrement(t *testing.T) {
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

	city := tc.ProjectID + "-DC"
	dc := client.Collection("cities").Doc(city)
	data := map[string]int{"population": 100}
	dc.Set(ctx, data)

	if err := updateDocumentIncrement(projectID, city); err != nil {
		t.Fatalf("updateDocumentIncrement: %v", err)
	}

	ref, err := dc.Get(ctx)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if err := ref.DataTo(&data); err != nil {
		t.Fatalf("DataTo: %v", err)
	}
	if got, want := data["population"], 150; got != want {
		t.Fatalf("updateDocumentIncrement DC population got %d, want %d", got, want)
	}
}
