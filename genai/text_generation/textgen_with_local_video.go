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

// [START googlegenaisdk_textgen_with_local_video]
import (
	"context"
	"fmt"
	"io"
	"os"

	genai "google.golang.org/genai"
)

// generateWithLocalVideo shows how to generate text using a local video input.
func generateWithLocalVideo(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// Read local video file content
	data, err := os.ReadFile("describe_video_content.mp4")
	if err != nil {
		return fmt.Errorf("failed to read local video: %w", err)
	}

	modelName := "gemini-2.5-flash"
	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: `Write a short and engaging blog post based on this video.`},
				{InlineData: &genai.Blob{
					MIMEType: "video/mp4",
					Data:     data,
				}},
			},
		},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()
	fmt.Fprintln(w, respText)

	// Example response:
	// Finding Your Flow: The Focused Ascent
	//
	// Ever watched someone scale an indoor climbing wall and been captivated by their precision and power? This video perfectly captures that intense focus and calculated movement.
	//
	// Our climber isn't just pulling himself up; he's engaging in a dynamic dance with gravity. Every reach, every foot placement, every clip of the rope is a deliberate part of solving the route's puzzle. You can almost feel the concentration as his eyes scan for the next optimal hold, his muscles working in unison to propel him upwards.
	//
	// Indoor climbing....
	// ...

	return nil
}

// [END googlegenaisdk_textgen_with_local_video]
