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

func TestCreateAndDeleteGlossary(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer
	var got string

	location := "us-central1"
	glossaryID := fmt.Sprintf("create_and_delete_glossary-%v", uuid.New().ID())
	glossaryInputURI := "gs://cloud-samples-data/translation/glossary_ja.csv"

	// Create a glossary.
	buf.Reset()
	if err := createGlossary(&buf, tc.ProjectID, location, glossaryID, glossaryInputURI); err != nil {
		t.Fatalf("createGlossary: %v", err)
	}
	got = buf.String()
	if !strings.Contains(got, "Created") {
		t.Fatalf("Got '%s', expected to contain 'Created'", got)
	}
	if !strings.Contains(got, glossaryID) {
		t.Fatalf("Got '%s', expected to contain '%s'", got, glossaryID)
	}
	if !strings.Contains(got, glossaryInputURI) {
		t.Fatalf("Got '%s', expected to contain '%s'", got, glossaryInputURI)
	}

	// Delete the glossary.
	buf.Reset()
	if err := deleteGlossary(&buf, tc.ProjectID, location, glossaryID); err != nil {
		t.Fatalf("deleteGlossary: %v", err)
	}
	got = buf.String()
	if !strings.Contains(got, "Deleted") {
		t.Fatalf("Got '%s', expected to contain 'Deleted'", got)
	}
	if !strings.Contains(got, location) {
		t.Fatalf("Got '%s', expected to contain '%s'", got, location)
	}
	if !strings.Contains(got, glossaryID) {
		t.Fatalf("Got '%s', expected to contain '%s'", got, glossaryID)
	}
}
