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
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestGenerateEmbeddings(t *testing.T) {
	tc := testutil.SystemTest(t)
	apiEndpoint := "us-central1-aiplatform.googleapis.com:443"
	model := "text-embedding-004"
	texts := []string{"banana muffins? ", "banana bread? banana muffins?"}
	dimensionality := 5
	embeddings, err := embedTextsPreview(apiEndpoint, tc.ProjectID, model, texts, "QUESTION_ANSWERING", &dimensionality)
	if err != nil {
		t.Fatal(err)
	}
	if len(embeddings) != len(texts) || len(embeddings[0]) != dimensionality {
		t.Errorf("len(embeddings), len(embeddings[0]) = %d, %d, want %d, %d", len(embeddings), len(embeddings[0]), len(texts), dimensionality)
	}
}

func TestGenerateEmbeddingsPreview(t *testing.T) {
	tc := testutil.SystemTest(t)
	apiEndpoint := "us-central1-aiplatform.googleapis.com:443"
	model := "text-embedding-preview-0409"
	texts := []string{"banana muffins? ", "banana bread? banana muffins?"}
	dimensionality := 5
	embeddings, err := embedTextsPreview(apiEndpoint, tc.ProjectID, model, texts, "QUESTION_ANSWERING", &dimensionality)
	if err != nil {
		t.Fatal(err)
	}
	if len(embeddings) != len(texts) || len(embeddings[0]) != dimensionality {
		t.Errorf("len(embeddings), len(embeddings[0]) = %d, %d, want %d, %d", len(embeddings), len(embeddings[0]), len(texts), dimensionality)
	}
}
