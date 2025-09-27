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

// Package image_generation shows how to use the GenAI SDK to generate images from prompt.

package image_generation

// [START googlegenaisdk_imggen_style_reference_with_txt_img]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateStyleRefWithText demonstrates how to generate an image using a style reference.
func generateStyleRefWithText(w io.Writer, outputGCSURI string) error {
	//outputGCSURI = "gs://your-bucket/your-prefix"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// Create a style reference image of a neon sign stored in Google Cloud Storage
	// using https://storage.googleapis.com/cloud-samples-data/generative-ai/image/neon.png
	styleRefImg := &genai.StyleReferenceImage{
		ReferenceID: 1,
		ReferenceImage: &genai.Image{
			GCSURI: "gs://cloud-samples-data/generative-ai/image/neon.png",
		},
		Config: &genai.StyleReferenceConfig{
			StyleDescription: "neon sign",
		},
	}

	// prompt that references the style image with [1]
	prompt := "generate an image of a neon sign [1] with the words: have a great day"
	modelName := "imagen-3.0-capability-001"

	// EditImage takes: ctx, model, prompt, referenceImages []ReferenceImage, config *EditImageConfig
	resp, err := client.Models.EditImage(ctx,
		modelName,
		prompt,
		[]genai.ReferenceImage{
			styleRefImg,
		},
		&genai.EditImageConfig{
			EditMode:          genai.EditModeDefault,
			NumberOfImages:    1,
			SafetyFilterLevel: genai.SafetyFilterLevelBlockMediumAndAbove,
			PersonGeneration:  genai.PersonGenerationAllowAdult,
			OutputGCSURI:      outputGCSURI,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to edit image: %w", err)
	}

	if len(resp.GeneratedImages) == 0 || resp.GeneratedImages[0].Image == nil {
		return fmt.Errorf("no generated images returned")
	}

	// Grab the first generated image URI and print it
	uri := resp.GeneratedImages[0].Image.GCSURI
	fmt.Fprintln(w, uri)

	// Example response:
	// gs://your-bucket/your-prefix
	return nil
}

// [END googlegenaisdk_imggen_style_reference_with_txt_img]
