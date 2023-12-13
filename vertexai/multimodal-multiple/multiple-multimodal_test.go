// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/vertexai/genai"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestGenerateMultimodalContent(t *testing.T) {
	t.Skip("TODO(muncus): remove skip")
	tc := testutil.SystemTest(t)

	projectID := tc.ProjectID
	location := "us-central1"

	modelName := "gemini-pro-vision"
	temperature := 0.8

	colosseum, err := partFromImageURL("https://storage.googleapis.com/cloud-samples-data/vertex-ai/llm/prompts/landmark1.png")
	if err != nil {
		t.Fatal(err)
	}
	// forbidden city
	forbiddenCity, err := partFromImageURL("https://storage.googleapis.com/cloud-samples-data/vertex-ai/llm/prompts/landmark2.png")
	if err != nil {
		t.Fatal(err)
	}
	// new image
	newImage, err := partFromImageURL("https://storage.googleapis.com/cloud-samples-data/vertex-ai/llm/prompts/landmark3.png")
	if err != nil {
		t.Fatal(err)
	}

	// create a multimodal (multipart) prompt
	prompt := []genai.Part{
		colosseum,
		genai.Text("city: Rome, Landmark: the Colosseum "),
		forbiddenCity,
		genai.Text("city: Beijing, Landmark: the Forbidden City "),
		newImage,
	}

	if projectID == "" {
		t.Fatal("require environment variable GOOGLE_CLOUD_PROJECT")
	}

	var buf bytes.Buffer
	if err := generateMultimodalContent(os.Stdout, prompt, projectID, location, modelName, float32(temperature)); err != nil {
		t.Fatal(err)
	}

	if got := buf.String(); !strings.Contains(got, "generated response") {
		t.Error("generated text content not found in response")
	}
}
