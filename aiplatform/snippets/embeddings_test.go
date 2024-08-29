// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package snippets

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestGenerateEmbeddings(t *testing.T) {
	tc := testutil.SystemTest(t)
	texts := []string{"banana muffins? ", "banana bread? banana muffins?"}
	dimensionality := 5
	location := "us-central1"
	var buf bytes.Buffer

	err := embedTexts(&buf, tc.ProjectID, location)
	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	if output != fmt.Sprintf("Dimensionality: %d. Embeddings length: %d", dimensionality, len(texts)) {
		t.Error("Embeddings length and dimensionality doesn't match")
	}
}

func TestGenerateEmbeddingsPreview(t *testing.T) {
	tc := testutil.SystemTest(t)
	texts := []string{"banana muffins? ", "banana bread? banana muffins?"}
	location := "us-central1"
	dimensionality := 5
	var buf bytes.Buffer

	err := embedTextsPreview(&buf, tc.ProjectID, location)
	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	if output != fmt.Sprintf("Dimensionality: %d. Embeddings length: %d", dimensionality, len(texts)) {
		t.Error("Embeddings length and dimensionality doesn't match")
	}
}
