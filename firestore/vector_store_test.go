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
	"testing"

	"cloud.google.com/go/firestore"
)

func testStoreVectors(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := storeVectors(buf, projectID); err != nil {
		t.Errorf("storeVectors: %v", err)
	}

	// clean up
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Errorf("firestore.NewClient: %v", err)
	}
	t.Cleanup(func() { client.Close() })
	collName := "coffee-beans"
	docs, err := client.Collection(collName).DocumentRefs(ctx).GetAll()
	if err != nil {
		t.Errorf("GetAll: %v", err)
	}
	for _, doc := range docs {
		if _, err := doc.Delete(ctx); err != nil {
			t.Errorf("Delete: %v", err)
		}
	}
}
