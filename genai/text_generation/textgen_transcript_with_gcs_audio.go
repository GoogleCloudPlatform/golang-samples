// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package text_generation shows examples of generating text using the GenAI SDK.
package text_generation

// [START googlegenaisdk_textgen_transcript_with_gcs_audio]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateAudioTranscript shows how to generate an audio transcript.
func generateAudioTranscript(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash"
	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: `Transcribe the interview, in the format of timecode, speaker, caption.
Use speaker A, speaker B, etc. to identify speakers.`},
			{FileData: &genai.FileData{
				FileURI:  "gs://cloud-samples-data/generative-ai/audio/pixel.mp3",
				MIMEType: "audio/mpeg",
			}},
		},
			Role: "user"},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)

	// Example response:
	// 00:00:00, A: your devices are getting better over time.
	// 00:01:13, A: And so we think about it across the entire portfolio from phones to watch, ...
	// ...

	return nil
}

// [END googlegenaisdk_textgen_transcript_with_gcs_audio]
