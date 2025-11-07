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

package image_generation

import (
	"bytes"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestImageGeneration(t *testing.T) {
	tc := testutil.SystemTest(t)

	t.Setenv("GOOGLE_GENAI_USE_VERTEXAI", "1")
	t.Setenv("GOOGLE_CLOUD_LOCATION", "global")
	t.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	buf := new(bytes.Buffer)

	t.Run("generate multimodal flash content with text and image", func(t *testing.T) {
		buf.Reset()
		err := generateMMFlashWithText(buf)
		if err != nil {
			t.Fatalf("generateMMFlashWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate mmflash text and image recipe", func(t *testing.T) {
		buf.Reset()
		err := generateMMFlashTxtImgWithText(buf)
		if err != nil {
			t.Fatalf("generateMMFlashTxtImgWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("style customization with style reference", func(t *testing.T) {
		buf.Reset()
		// TODO(developer): update with your bucket
		outputGCSURI := "gs://your-bucket/your-prefix"

		err := generateStyleRefWithText(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateStyleRefWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected printed output, got empty")
		}
	})

	t.Run("canny edge customization with text+image", func(t *testing.T) {
		buf.Reset()
		// TODO(developer): update with your bucket
		outputGCSURI := "gs://your-bucket/your-prefix"

		err := generateCannyCtrlTypeWithText(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateCannyCtrlTypeWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate image with scribble control type", func(t *testing.T) {
		buf.Reset()
		// TODO(developer): update with your bucket
		outputGCSURI := "gs://your-bucket/your-prefix"

		err := generateScribbleCtrlTypeWithText(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateScribbleCtrlTypeWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("subject customization with control reference", func(t *testing.T) {
		buf.Reset()
		// TODO(developer): update with your bucket
		outputGCSURI := "gs://your-bucket/your-prefix"

		err := generateSubjRefCtrlReferWithText(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateSubjRefCtrlReferWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate style transfer customization with raw reference", func(t *testing.T) {
		buf.Reset()
		// TODO(developer): update with your bucket
		outputGCSURI := "gs://your-bucket/your-prefix"

		err := generateRawReferWithText(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateRawReferWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})
}
