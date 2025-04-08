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

	location := "us-central1"
	glossaryID := fmt.Sprintf("create_and_delete_glossary-%v", uuid.New().ID())
	glossaryInputURI := "gs://cloud-samples-data/translation/glossary_ja.csv"

	// Create a glossary.
	var buf bytes.Buffer
	if err := createGlossary(&buf, tc.ProjectID, location, glossaryID, glossaryInputURI); err != nil {
		t.Fatalf("createGlossary: %v", err)
	}
	if got, want := buf.String(), "Created"; !strings.Contains(got, want) {
		t.Errorf("createGlossary got:\n----\n%s----\nWant to contain:\n----\n%s\n----", got, want)
	}
	if got, want := buf.String(), glossaryID; !strings.Contains(got, want) {
		t.Errorf("createGlossary got:\n----\n%s----\nWant to contain:\n----\n%s\n----", got, want)
	}
	if got, want := buf.String(), glossaryInputURI; !strings.Contains(got, want) {
		t.Errorf("createGlossary got:\n----\n%s----\nWant to contain:\n----\n%s\n----", got, want)
	}

	// Delete the glossary.
	buf.Reset()
	if err := deleteGlossary(&buf, tc.ProjectID, location, glossaryID); err != nil {
		t.Fatalf("deleteGlossary: %v", err)
	}
	if got, want := buf.String(), "Deleted"; !strings.Contains(got, want) {
		t.Errorf("deleteGlossary got:\n----\n%s----\nWant to contain:\n----\n%s\n----", got, want)
	}
	if got, want := buf.String(), location; !strings.Contains(got, want) {
		t.Errorf("deleteGlossary got:\n----\n%s----\nWant to contain:\n----\n%s\n----", got, want)
	}
	if got, want := buf.String(), glossaryID; !strings.Contains(got, want) {
		t.Errorf("deleteGlossary got:\n----\n%s----\nWant to contain:\n----\n%s\n----", got, want)
	}
}
