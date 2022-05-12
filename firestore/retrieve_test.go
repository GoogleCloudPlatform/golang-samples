// Copyright 2021 Google LLC
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
	"time"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestRetrieve(t *testing.T) {
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

	if err = prepareRetrieve(ctx, client); err != nil {
		t.Fatalf("prepareRetrieve: %v", err)
	}

	testutil.Retry(t, 3, 5*time.Second, func(r *testutil.R) {
		_, err = docAsMap(ctx, client)
		if err != nil {
			r.Errorf("Cannot get doc as map: %v", err)
		}
	})

	_, err = docAsEntity(ctx, client)
	if err != nil {
		t.Fatalf("Cannot get doc as entity: %v", err)
	}

	if err = multipleDocs(ctx, client); err != nil {
		t.Fatalf("multipleDocs: %v", err)
	}
	if err = allDocs(ctx, client); err != nil {
		t.Fatalf("allDocs: %v", err)
	}
	if err = getCollections(ctx, client); err != nil {
		t.Fatalf("getCollections: %v", err)
	}
}
