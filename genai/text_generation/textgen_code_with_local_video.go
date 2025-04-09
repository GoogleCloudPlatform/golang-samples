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

// [START googlegenaisdk_textgen_code_with_local_video]
import (
	"context"
	"fmt"
	"io"
	"os"

	genai "google.golang.org/genai"
)

// generateWithLocalVideo shows how to generate text using a local video file as input.
func generateWithLocalVideo(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// TODO(Developer): Update the path to file (video source:
	//   https://storage.googleapis.com/cloud-samples-data/generative-ai/video/describe_video_content.mp4)
	videoBytes, err := os.ReadFile("./describe_video_content.mp4")
	if err != nil {
		return fmt.Errorf("failed to read video file: %w", err)
	}

	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: "Write a short and engaging blog post based on this video."},
			{InlineData: &genai.Blob{
				Data:     videoBytes,
				MIMEType: "video/mp4",
			}},
		}},
	}
	modelName := "gemini-2.0-flash-001"

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText, err := resp.Text()
	if err != nil {
		return fmt.Errorf("failed to convert model response to text: %w", err)
	}
	fmt.Fprintln(w, respText)

	// Example response:
	// Okay, here's a blog post inspired by the provided image:
	//
	// **Title: Conquer Your Fears, Conquer the Wall!**
	// ...

	return nil
}

// [END googlegenaisdk_textgen_code_with_local_video]
