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

// [START googlegenaisdk_imggen_mmflash_txt_and_img_with_txt]
import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"google.golang.org/genai"
)

// generateMMFlashTxtImgWithText demonstrates how to generate an illustrated recipe
// combining text and image outputs into a markdown file.
func generateMMFlashTxtImgWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash-image-preview"
	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{Text: "Generate an illustrated recipe for a paella. " +
					"Create images to go alongside the text as you generate the recipe."},
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
			CandidateCount: int32(1),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return fmt.Errorf("no candidates returned")
	}

	outputFolder := "testdata"

	// Create the markdown file
	mdFile := filepath.Join(outputFolder, "paella-recipe.md")
	fp, err := os.Create(mdFile)
	if err != nil {
		return fmt.Errorf("failed to create markdown file: %w", err)
	}
	defer fp.Close()

	for i, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			if _, err := fp.WriteString(part.Text); err != nil {
				return fmt.Errorf("failed to write text: %w", err)
			}
		} else if part.InlineData != nil {
			imgFile := filepath.Join(outputFolder, fmt.Sprintf("example-image-%d.png", i+1))
			if err := os.WriteFile(imgFile, part.InlineData.Data, 0644); err != nil {
				return fmt.Errorf("failed to save image: %w", err)
			}
			if _, err := fp.WriteString(fmt.Sprintf("![image](%s)", filepath.Base(imgFile))); err != nil {
				return fmt.Errorf("failed to write image reference: %w", err)
			}
		}
	}

	fmt.Fprintln(w, mdFile)

	// Example response:
	//  A markdown page for a Paella recipe (`paella-recipe.md`) has been generated.
	//  It includes detailed steps and several images illustrating the cooking process.
	return nil
}

// [END googlegenaisdk_imggen_mmflash_txt_and_img_with_txt]
