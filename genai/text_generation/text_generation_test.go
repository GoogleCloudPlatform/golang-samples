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

package text_generation

import (
	"bytes"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestTextGeneration(t *testing.T) {
	tc := testutil.SystemTest(t)

	t.Setenv("GOOGLE_GENAI_USE_VERTEXAI", "1")
	t.Setenv("GOOGLE_CLOUD_LOCATION", "us-central1")
	t.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	buf := new(bytes.Buffer)

	t.Run("generate text with text prompt", func(t *testing.T) {
		buf.Reset()
		err := generateWithText(buf)
		if err != nil {
			t.Fatalf("generateWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate text with text and image prompt", func(t *testing.T) {
		buf.Reset()
		err := generateWithTextImage(buf)
		if err != nil {
			t.Fatalf("generateWithTextImage failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate text with text prompt and custom configuration", func(t *testing.T) {
		buf.Reset()
		err := generateWithConfig(buf)
		if err != nil {
			t.Fatalf("generateWithConfig failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate text with text prompt and system instructions", func(t *testing.T) {
		buf.Reset()
		err := generateWithSystem(buf)
		if err != nil {
			t.Fatalf("generateWithSystem failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate text stream with text prompt", func(t *testing.T) {
		buf.Reset()
		err := generateWithTextStream(buf)
		if err != nil {
			t.Fatalf("generateWithTextStream failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})
}
