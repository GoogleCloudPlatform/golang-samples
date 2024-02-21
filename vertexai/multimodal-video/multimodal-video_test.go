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
	"context"
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

	ctx := context.Background()

	// create video part
	video, err := partFromGCSURI(ctx, "gs://cloud-samples-data/video/animals.mp4")
	if err != nil {
		t.Fatalf("unable to process media: %v", err)
	}
	prompt := []genai.Part{
		genai.Text("What is in the video? "),
		video,
	}

	if projectID == "" {
		t.Fatal("require environment variable GOOGLE_CLOUD_PROJECT")
	}

	var buf bytes.Buffer
	if err := generateMultimodalContent(ctx, os.Stdout, prompt, projectID, location, modelName); err != nil {
		t.Fatal(err)
	}

	if got := buf.String(); !strings.Contains(got, "generated response") {
		t.Error("generated text content not found in response")
	}
}
