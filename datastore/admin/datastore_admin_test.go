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

package samples

import (
	"context"
	"io"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestAdmin(t *testing.T) {
	// Roles to be set in your Service Account and App Engine default service account
	// to run this test:
	// `Datastore Import Export Admin`, or `Cloud Datastore Owner`, or `Owner`,
	// `Storage Admin`, or `Owner`.
	// See https://cloud.google.com/datastore/docs/export-import-entities#permissions for full details
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := clientCreate(io.Discard)
	if err != nil {
		t.Fatalf("clientCreate: %v", err)
	}
	defer client.Close()

	indices, err := indexList(io.Discard, tc.ProjectID)
	if err != nil {
		t.Fatalf("indexList: %v", err)
	}
	if len(indices) == 0 {
		t.Skip("Skipping datastore test. At least one index should present in database.")
	}
	want := indices[0].IndexId
	// Get the first index from the list of indexes.
	got, err := indexGet(io.Discard, tc.ProjectID, want)
	if err != nil {
		t.Fatalf("indexGet: %v", err)
	}
	if got.IndexId != want {
		t.Fatalf("Unexpected indexID: got %v, want %v", got.IndexId, want)
	}
	// Create bucket for Export/Import entities.
	bucketName := testutil.TestBucket(ctx, t, tc.ProjectID, "datastore-admin")

	resp, err := entitiesExport(io.Discard, tc.ProjectID, "gs://"+bucketName)
	if err != nil {
		t.Fatalf("entitiesExport: %v", err)
	}
	if err := entitiesImport(io.Discard, tc.ProjectID, resp.OutputUrl); err != nil {
		t.Fatalf("entitiesImport: %v", err)
	}
}
