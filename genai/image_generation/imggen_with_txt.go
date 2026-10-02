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

// [START googlegenaisdk_imggen_with_txt]
import (
	"context"
	"fmt"
	"io"
	"os"

	"google.golang.org/genai"
)

// generateImageWithText demonstrates how to generate an image from a text prompt.
func generateImageWithText(w io.Writer) error {
	// TODO(developer): Update below line
	outputFile := "dog_newspaper.png"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash-image"
	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{Text: "A dog reading a newspaper"},
			},
			Role: genai.RoleUser,
		},
	}

	resp, err := client.Models.GenerateContent(ctx,
		modelName,
		contents,
		&genai.GenerateContentConfig{
			ResponseModalities: []string{
				string(genai.ModalityImage),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("no image generated")
	}

	part := resp.Candidates[0].Content.Parts[0]
	if part.InlineData == nil {
		return fmt.Errorf("no image generated")
	}

	if err := os.WriteFile(outputFile, part.InlineData.Data, 0644); err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	fmt.Fprintln(w, len(part.InlineData.Data))

	// Example response:
	// Created output image using 6098201 bytes
	return nil
}

// [END googlegenaisdk_imggen_with_txt]
