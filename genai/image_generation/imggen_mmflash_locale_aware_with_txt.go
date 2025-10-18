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

// [START googlegenaisdk_imggen_mmflash_locale_aware_with_txt]
import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"google.golang.org/genai"
)

// generateMMFlashLocaleAwareWithText demonstrates how to generate an image with locale awareness.
func generateMMFlashLocaleAwareWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash-image"
	prompt := "Generate a photo of a breakfast meal."
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

	outputPath := filepath.Join("testdata", "example-breakfast-meal.png")

	for _, part := range resp.Candidates[0].Content.Parts {
		switch {
		case part.Text != "":
			fmt.Fprintln(w, part.Text)
		case part.InlineData != nil:
			if err := os.WriteFile(outputPath, part.InlineData.Data, 0o644); err != nil {
				return fmt.Errorf("failed to save generated image: %w", err)
			}
		}
	}
	fmt.Fprintln(w, outputPath)

	// Example response:
	//  Here is a photo of a delicious breakfast meal for you! ...

	return nil
}

// [END googlegenaisdk_imggen_mmflash_locale_aware_with_txt]
