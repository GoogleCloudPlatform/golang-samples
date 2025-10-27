// Copyright 2025 Google LLC
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

package video_generation

import (
	"bytes"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestVideoGeneration(t *testing.T) {
	tc := testutil.SystemTest(t)

	t.Setenv("GOOGLE_GENAI_USE_VERTEXAI", "1")
	t.Setenv("GOOGLE_CLOUD_LOCATION", "us-central1")
	t.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	buf := new(bytes.Buffer)

	gcsOutputBucket := "HERE-go-bucket-samples-tests"
	prefix := "go_videogen_test/" + time.Now().Format("20060102-150405")
	outputGCSURI := "gs://" + gcsOutputBucket + "/" + prefix

	t.Run("generate video content with img", func(t *testing.T) {
		buf.Reset()
		err := generateVideoFromImage(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateVideoFromImage failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate video content with text", func(t *testing.T) {
		buf.Reset()
		err := generateVideoWithText(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateVideoWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

}
