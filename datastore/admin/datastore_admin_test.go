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
	"io/ioutil"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreate(t *testing.T) {
	client, err := clientCreate(ioutil.Discard)
	if err != nil {
		t.Fatalf("clientCreate: %v", err)
	}
	defer client.Close()

	tc := testutil.SystemTest(t)
	indices, err := indexList(ioutil.Discard, tc.ProjectID)
	if err != nil {
		t.Fatalf("indexList: %v", err)
	}
	if len(indices) == 0 {
		t.Skip("Skipping datastore test. At least one index should present in database.")
	}
	want := indices[0].IndexId
	// Get the first index from the list of indexes.
	got, err := indexGet(ioutil.Discard, tc.ProjectID, want)
	if err != nil {
		t.Fatalf("indexGet: %v", err)
	}
	if got.IndexId != want {
		t.Fatalf("Unexpected indexID: got %v, want %v", got.IndexId, want)
	}
}
