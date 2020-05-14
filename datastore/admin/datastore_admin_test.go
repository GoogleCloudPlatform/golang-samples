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
}

func TestIndexList(t *testing.T) {
	tc := testutil.SystemTest(t)
	if _, err := indexList(ioutil.Discard, tc.ProjectID); err != nil {
		t.Fatalf("indexList: %v", err)
	}
}

func TestIndexGet(t *testing.T) {
	tc := testutil.SystemTest(t)
	// Get first index from the list of indexes.
	indexes, _ := indexList(ioutil.Discard, tc.ProjectID)
	indexID := indexes[0].IndexId

	if err := indexGet(ioutil.Discard, tc.ProjectID, indexID); err != nil {
		t.Fatalf("indexGet: %v", err)
	}
}
