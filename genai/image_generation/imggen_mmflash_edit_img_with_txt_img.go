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

// Package image_generation shows how to use the GenAI SDK to generate images and text.
package image_generation

// [START googlegenaisdk_imggen_mmflash_edit_img_with_txt_img]
import (
	"context"
	"fmt"
	"io"
	"os"

	"google.golang.org/genai"
)

// generateImageMMFlashEditWithTextImg demonstrates editing an image with text and image inputs.
func generateImageMMFlashEditWithTextImg(w io.Writer) error {
	// TODO(developer): Update below lines
	outputFile := "testdata/bw-example-image.png"
	inputFile := "testdata/example-image-eiffel-tower.png"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	image, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read image: %w", err)
	}

	modelName := "gemini-2.5-flash-image-preview"
	prompt := "Edit this image to make it look like a cartoon."
	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: prompt},
				{InlineData: &genai.Blob{
					MIMEType: "image/png",
					Data:     image,
				}},
			},
		},
	}
	resp, err := client.Models.GenerateContent(ctx,
		modelName,
		contents,
		&genai.GenerateContentConfig{
			ResponseModalities: []string{
				string(genai.ModalityText),
				string(genai.ModalityImage),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return fmt.Errorf("no candidates returned")
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			fmt.Fprintln(w, part.Text)
		} else if part.InlineData != nil {
			if len(part.InlineData.Data) > 0 {
				if err := os.WriteFile(outputFile, part.InlineData.Data, 0644); err != nil {
					return fmt.Errorf("failed to save image: %w", err)
				}
				fmt.Fprintln(w, outputFile)
			}
		}
	}

	// Example response:
	// Here's the image of the Eiffel Tower and fireworks, cartoonized for you!
	// Cartoon-style edit:
	//  - Simplified the Eiffel Tower with bolder lines and slightly exaggerated proportions.
	//  - Brightened and saturated the colors of the sky, fireworks, and foliage for a more vibrant, cartoonish look.
	//  ....
	return nil
}

// [END googlegenaisdk_imggen_mmflash_edit_img_with_txt_img]
