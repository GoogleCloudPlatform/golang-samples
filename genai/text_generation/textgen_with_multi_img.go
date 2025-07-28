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

// [START googlegenaisdk_textgen_with_multi_img]
import (
	"context"
	"fmt"
	"io"
	"os"

	genai "google.golang.org/genai"
)

// generateWithMultiImg shows how to generate text using multiple image inputs.
func generateWithMultiImg(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// TODO(Developer): Update the path to file (image source:
	//   https://storage.googleapis.com/cloud-samples-data/generative-ai/image/latte.jpg )
	imageBytes, err := os.ReadFile("./latte.jpg")
	if err != nil {
		return fmt.Errorf("failed to read image: %w", err)
	}

	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: "Write an advertising jingle based on the items in both images."},
			{FileData: &genai.FileData{
				// Image source: https://storage.googleapis.com/cloud-samples-data/generative-ai/image/scones.jpg
				FileURI:  "gs://cloud-samples-data/generative-ai/image/scones.jpg",
				MIMEType: "image/jpeg",
			}},
			{InlineData: &genai.Blob{
				Data:     imageBytes,
				MIMEType: "image/jpeg",
			}},
		}},
	}
	modelName := "gemini-2.5-flash"

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)

	// Example response:
	// Okay, here's an advertising jingle inspired by the blueberry scones, coffee, flowers, chocolate cake, and latte:
	//
	// (Upbeat, jazzy music)
	// ...

	return nil
}

// [END googlegenaisdk_textgen_with_multi_img]
