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

// generateWithMultiLocalImages shows how to generate text using multiple local image inputs.
func generateWithMultiLocalImages(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// Read local image files
	image1, err := os.ReadFile("test_data/latte.jpg")
	if err != nil {
		return fmt.Errorf("failed to read image1: %w", err)
	}
	image2, err := os.ReadFile("test_data/scones.jpg")
	if err != nil {
		return fmt.Errorf("failed to read image2: %w", err)
	}

	modelName := "gemini-2.5-flash"
	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: "Generate a list of all the objects contained in both images."},
				{InlineData: &genai.Blob{
					MIMEType: "image/jpeg",
					Data:     image1,
				}},
				{InlineData: &genai.Blob{
					MIMEType: "image/jpeg",
					Data:     image2,
				}},
			},
		},
	}

	// Call the model
	resp, err := client.Models.GenerateContent(ctx, modelName, contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	fmt.Fprintln(w, resp.Text())

	// Example response:
	// Here is a list of all the distinct objects found in both images:
	// 1.  **Coffee** (in mugs/cups; one is clearly a latte with heart art, others are also coffee/latte)
	// 2.  **Mug(s)/Cup(s)** (yellow in the top image, white in the bottom image)
	// 3.  **Cake** (sliced, in the top image)
	// 4.  **Plate** (white, under the cake slice in the top image)
	// 5.  **Fork** (partially visible on the plate in the top image)
	// 6.  **Scones/Biscuits** (blueberry, in the bottom image)
	// 7.  **Blueberries** (scattered and in a bowl in the bottom image)
	// 8.  **Bowl** (small, dark, holding blueberries in the bottom image)
	// 9.  **Spoon** (silver, with "LET'S JAM" inscription, in the bottom image)
	// 10. **Flowers** (peonies, in the bottom image)
	// 11. **Leaves** (green, possibly mint, in the bottom image)
	// 12. **Paper** (parchment or wax paper, in the bottom image)
	// 13. **Table/Surface** (wooden in the top image, textured/painted in the bottom image)
	// ...

	return nil
}

// [END googlegenaisdk_textgen_with_multi_local_img]
