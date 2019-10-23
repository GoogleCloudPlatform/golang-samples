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

package v3

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

func TestListGlossaries(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer
	var got string

	location := "us-central1"
	glossaryID := fmt.Sprintf("get_glossary-%v", uuid.New().ID())
	glossaryInputURI := "gs://cloud-samples-data/translation/glossary_ja.csv"

	// Create a glossary.
	if err := createGlossary(&buf, tc.ProjectID, location, glossaryID, glossaryInputURI); err != nil {
		t.Fatalf("createGlossary: %v", err)
	}
	defer deleteGlossary(&buf, tc.ProjectID, location, glossaryID)

	// Check the glossaries.
	if err := listGlossaries(&buf, tc.ProjectID, location); err != nil {
		t.Fatalf("listGlossaries: %v", err)
	}
	got = buf.String()
	if !strings.Contains(got, glossaryID) {
		t.Fatalf("Got '%s', expected to contain '%s'", got, glossaryID)
	}
	if !strings.Contains(got, glossaryInputURI) {
		t.Fatalf("Got '%s', expected to contain '%s'", got, glossaryInputURI)
	}
}
