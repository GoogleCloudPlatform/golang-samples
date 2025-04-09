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

// [START googlegenaisdk_textgen_with_multi_local_img]
import (
	"context"
	"fmt"
	"io"
	"os"

	genai "google.golang.org/genai"
)

// generateWithMultiLocalImg shows how to generate text using multiple local images as inputs.
func generateWithMultiLocalImg(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// TODO(Developer): Update with paths to your image files.
	imagePath := filepath.Join(getMedia(), "latte.jpg")
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}
	imagePath := filepath.Join(getMedia(), "scones.jpg")
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: "Write an advertising jingle based on the items in both images."},
			{InlineData: &genai.Blob{
				Data:     imageBytes1,
				MIMEType: "image/jpeg",
			}},
			{InlineData: &genai.Blob{
				Data:     imageBytes2,
				MIMEType: "image/jpeg",
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
	// Okay, here's a jingle inspired by the images of cake, coffee, and blueberry scones:
	//
	// (Upbeat, folksy music)
	// ...

	return nil
}

// [END googlegenaisdk_textgen_with_multi_local_img]
