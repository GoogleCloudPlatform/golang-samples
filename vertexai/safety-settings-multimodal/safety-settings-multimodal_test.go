// Copyright 2023 Google LLC
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

package main

import (
	"bytes"
	"strings"
	"testing"

	"cloud.google.com/go/vertexai/genai"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestGenerateContent(t *testing.T) {
	t.Skip("TODO(muncus): remove skip")
	tc := testutil.SystemTest(t)

	projectID := tc.ProjectID
	location := "us-central1"

	model := "gemini-pro-vision"
	temp := 0.8

	cat, _ := partFromImageURL("https://storage.googleapis.com/cloud-samples-data/generative-ai/image/320px-Felis_catus-cat_on_snow.jpg")

	// create a multipart (multimodal) prompt
	prompt := []genai.Part{
		genai.Text("say something nice about this "),
		cat,
	}

	if projectID == "" {
		t.Fatal("require environment variable GOOGLE_CLOUD_PROJECT")
	}

	var buf bytes.Buffer
	if err := generateContent(&buf, prompt, projectID, location, model, float32(temp)); err != nil {
		t.Fatal(err)
	}

	if got := buf.String(); !strings.Contains(got, "generate-content response") {
		t.Error("generated text content not found in response")
	}
}
