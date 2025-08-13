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

	t.Run("generate with text prompt", func(t *testing.T) {
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

	t.Run("generate with text prompt and custom configuration", func(t *testing.T) {
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

	t.Run("generate with text prompt and system instructions", func(t *testing.T) {
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

	t.Run("generate stream with text prompt", func(t *testing.T) {
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

	t.Run("generate with text and image prompt", func(t *testing.T) {
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

	t.Run("generate with pdf file input", func(t *testing.T) {
		buf.Reset()
		err := generateWithPDF(buf)
		if err != nil {
			t.Fatalf("generateWithPDF failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate with video file input (no sound)", func(t *testing.T) {
		buf.Reset()
		err := generateWithMuteVideo(buf)
		if err != nil {
			t.Fatalf("generateWithMuteVideo failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate with video file input", func(t *testing.T) {
		buf.Reset()
		err := generateWithVideo(buf)
		if err != nil {
			t.Fatalf("generateWithVideo failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate with audio file input", func(t *testing.T) {
		buf.Reset()
		err := generateWithAudio(buf)
		if err != nil {
			t.Fatalf("generateWithAudio failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate an audio transcript", func(t *testing.T) {
		buf.Reset()
		err := generateAudioTranscript(buf)
		if err != nil {
			t.Fatalf("generateAudioTranscript failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate with YT video file input", func(t *testing.T) {
		buf.Reset()
		err := generateWithYTVideo(buf)
		if err != nil {
			t.Fatalf("generateWithYTVideo failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate with routing", func(t *testing.T) {
		t.Skip("skipping because of model used in this test. The model 'model-optimizer-exp-04-09' is not consistently available in all test environments.")
		buf.Reset()
		err := generateWithRouting(buf)
		if err != nil {
			t.Fatalf("generateWithRouting failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate chat stream with text prompt", func(t *testing.T) {
		buf.Reset()
		err := generateChatStreamWithText(buf)
		if err != nil {
			t.Fatalf("generateChatStreamWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate Text With PDF", func(t *testing.T) {
		buf.Reset()
		err := generateTextWithPDF(buf)
		if err != nil {
			t.Fatalf("generateTextWithPDF failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate thinking with text prompt", func(t *testing.T) {
		buf.Reset()
		err := generateThinkingWithText(buf)
		if err != nil {
			t.Fatalf("generateThinkingWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate with model optimizer", func(t *testing.T) {
		buf.Reset()
		err := generateModelOptimizerWithTxt(buf)
		if err != nil {
			t.Fatalf("generateModelOptimizerWithTxt failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

}
