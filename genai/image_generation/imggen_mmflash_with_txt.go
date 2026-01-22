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

// [START googlegenaisdk_imggen_mmflash_with_txt]
import (
	"context"
	"fmt"
	"io"
	"os"

	"google.golang.org/genai"
)

// generateMMFlashWithText demonstrates how to generate both text and image outputs.
func generateMMFlashWithText(w io.Writer) error {
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
				{Text: "Generate an image of the Eiffel tower with fireworks in the background."},
			},
			Role: genai.RoleUser,
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
			SafetySettings: []*genai.SafetySetting{
				{Method: genai.HarmBlockMethodProbability},
				{Category: genai.HarmCategoryDangerousContent},
				{Threshold: genai.HarmBlockThresholdBlockMediumAndAbove},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return fmt.Errorf("no candidates returned")
	}
	var fileName string
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			fmt.Fprintln(w, part.Text)
		} else if part.InlineData != nil {
			fileName = "example-image-eiffel-tower.png"
			if err := os.WriteFile(fileName, part.InlineData.Data, 0o644); err != nil {
				return fmt.Errorf("failed to save image: %w", err)
			}
		}
	}
	fmt.Fprintln(w, fileName)

	// Example response:
	// I will generate an image of the Eiffel Tower at night, with a vibrant display of
	// colorful fireworks exploding in the dark sky behind it.
	// ....
	return nil
}

// [END googlegenaisdk_imggen_mmflash_with_txt]
