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

func TestDistributedCounter(t *testing.T) {
	testutil.EndToEndTest(t)
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

	docRef := client.Collection("counter_samples").Doc("DCounter")

	dc := Counter{3}
	dc.initCounter(ctx, docRef)

	shards := docRef.Collection("shards")
	docRefs, err := shards.DocumentRefs(ctx).GetAll()
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if l := len(docRefs); l != 3 {
		t.Fatalf("created %d shards, want 3", l)
	}

	dc.incrementCounter(ctx, docRef)
	dc.incrementCounter(ctx, docRef)

	total, err := dc.getCount(ctx, docRef)
	if err != nil {
		t.Fatalf("getCount: %v", err)
	}
	if total != 2 {
		t.Fatalf("got total = %d, want 2", total)
	}
}
