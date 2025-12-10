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

	t.Run("generate image content with text", func(t *testing.T) {
		buf.Reset()
		err := generateImageWithText(buf)
		if err != nil {
			t.Fatalf("generateImageWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate mmflash image content with text and image", func(t *testing.T) {
		buf.Reset()
		err := generateImageMMFlashEditWithTextImg(buf)
		if err != nil {
			t.Fatalf("generateImageMMFlashEditWithTextImg failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})
}
