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

// Package count_tokens shows examples of counting tokens using the GenAI SDK.
package count_tokens

// TODO[vburlaka@google.com,msampathkumar@google.com]: Remove `count_tokens` region tags after Feb 2025
// [START googlegenaisdk_count_tokens_with_txt_img_vid]
// [START googlegenaisdk_counttoken_with_txt_vid]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// countWithTxtAndVid shows how to count tokens with text and video inputs.
func countWithTxtAndVid(w io.Writer) error {
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
			{Text: "Provide a description of the video."},
			{FileData: &genai.FileData{
				FileURI:  "gs://cloud-samples-data/generative-ai/video/pixel8.mp4",
				MIMEType: "video/mp4",
			}},
		}},
	}

	resp, err := client.Models.CountTokens(ctx, modelName, contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	fmt.Fprintf(w, "Total: %d\nCached: %d\n", resp.TotalTokens, resp.CachedContentTokenCount)

	// Example response:
	// Total: 16252
	// Cached: 0

	return nil
}

// [END googlegenaisdk_counttoken_with_txt_vid]
// [END googlegenaisdk_count_tokens_with_txt_img_vid]
