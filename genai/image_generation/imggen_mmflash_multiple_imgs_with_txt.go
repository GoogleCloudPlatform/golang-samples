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

// [START googlegenaisdk_imggen_mmflash_multiple_imgs_with_txt]
import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"google.golang.org/genai"
)

// generateMMFlashMultipleImgsWithText demonstrates how to generate multiple images with text.
func generateMMFlashMultipleImgsWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash-image"
	prompt := "Generate 3 images a cat sitting on a chair."
	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{Text: prompt},
			},
			Role: "user",
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
		return fmt.Errorf("generate content failed: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return fmt.Errorf("no content was generated")
	}

	outputDir := filepath.Join("")

	imageCounter := 1
	for _, part := range resp.Candidates[0].Content.Parts {
		switch {
		case part.Text != "":
			fmt.Fprintln(w, part.Text)
		case part.InlineData != nil:
			filename := filepath.Join(outputDir, fmt.Sprintf("example-cats-0%d.png", imageCounter))
			if err := os.WriteFile(filename, part.InlineData.Data, 0o644); err != nil {
				return fmt.Errorf("failed to save generated image: %w", err)
			}
			fmt.Fprintln(w, filename)
			imageCounter++
		}
	}

	// Example response:
	//  Image 1: A fluffy calico cat with striking green eyes is perched elegantly on a vintage wooden
	//  chair with a woven seat. Sunlight streams through a nearby window, casting soft shadows and
	//  highlighting the cat's fur.
	//    #
	//  Image 2: A sleek black cat with intense yellow eyes is sitting upright on a modern, minimalist
	//  white chair. The background is a plain grey wall, putting the focus entirely on the feline's
	//  graceful posture.
	//    #
	//  Image 3: A ginger tabby cat with playful amber eyes is comfortably curled up asleep on a plush,
	//  oversized armchair upholstered in a soft, floral fabric. A corner of a cozy living room with a
	//  warm lamp in the background can be seen.

	return nil
}

// [END googlegenaisdk_imggen_mmflash_multiple_imgs_with_txt]
