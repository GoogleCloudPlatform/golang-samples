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

// [START googlegenaisdk_textgen_with_gcs_audio]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithAudio shows how to generate text using an audio input.
func generateWithAudio(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.0-flash-001"
	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: `Provide the summary of the audio file.
Summarize the main points of the audio concisely.
Create a chapter breakdown with timestamps for key sections or topics discussed.`},
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
	// Here is a summary and chapter breakdown of the audio file:
	//
	// **Summary:**
	//
	// The audio file is a "Made by Google" podcast episode discussing the Pixel Feature Drops, ...
	//
	// **Chapter Breakdown:**
	//
	// *   **0:00 - 0:54:** Introduction to the podcast and guests, Aisha Sharif and DeCarlos Love.
	// ...

	return nil
}

// [END googlegenaisdk_textgen_with_gcs_audio]
