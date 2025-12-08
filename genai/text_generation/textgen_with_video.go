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

// [START googlegenaisdk_textgen_with_video]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithVideo shows how to generate text using a video input.
func generateWithVideo(w io.Writer) error {
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
			{Text: `Analyze the provided video file, including its audio.
Summarize the main points of the video concisely.
Create a chapter breakdown with timestamps for key sections or topics discussed.`},
			{FileData: &genai.FileData{
				FileURI:  "gs://cloud-samples-data/generative-ai/video/pixel8.mp4",
				MIMEType: "video/mp4",
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
	// Here's an analysis of the provided video file:
	//
	// **Summary**
	//
	// The video features Saeka Shimada, a photographer in Tokyo, who uses the new Pixel phone ...
	//
	// **Chapter Breakdown**
	//
	// *   **0:00-0:05**: Introduction to Saeka Shimada and her work as a photographer in Tokyo.
	// ...

	return nil
}

// [END googlegenaisdk_textgen_with_video]
