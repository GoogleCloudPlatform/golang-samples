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

	"cloud.google.com/go/firestore"
)

func TestQuery(t *testing.T) {
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

	if err = prepareQuery(ctx, client); err != nil {
		t.Fatalf("prepareQuery: %v", err)
	}
	if err = paginateCursor(ctx, client); err != nil {
		t.Fatalf("paginateCursor: %v", err)
	}
	if err = createInQuery(ctx, client); err != nil {
		t.Fatalf("createInQuery: %v", err)
	}
	if err = createInQueryWithArray(ctx, client); err != nil {
		t.Fatalf("createInQueryWithArray: %v", err)
	}
	if err = createArrayContainsQuery(ctx, client); err != nil {
		t.Fatalf("createArrayContainsQuery: %v", err)
	}
	if err = createArrayContainsAnyQuery(ctx, client); err != nil {
		t.Fatalf("createArrayContainsAnyQuery: %v", err)
	}
	if err = createStartAtDocSnapshotQuery(ctx, client); err != nil {
		t.Fatalf("createStartAtDocSnapshotQuery: %v", err)
	}
}
