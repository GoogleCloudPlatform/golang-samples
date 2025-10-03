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

// [START googlegenaisdk_imggen_virtual_try_on_with_txt_img]
import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"google.golang.org/genai"
)

// generateImgVirtualTryOnWithTextImg demonstrates how to apply a product image to a person image.
func generateImgVirtualTryOnWithTextImg(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "virtual-try-on-preview-08-04"

	// Load local person image
	personBytes, err := os.ReadFile("testdata/man.png")
	if err != nil {
		return fmt.Errorf("failed to read person image: %w", err)
	}
	personImage := &genai.Image{
		ImageBytes: personBytes,
		MIMEType:   "image/png",
	}

	// Load local product image
	productBytes, err := os.ReadFile("testdata/sweater.jpg")
	if err != nil {
		return fmt.Errorf("failed to read product image: %w", err)
	}
	productImage := &genai.ProductImage{
		ProductImage: &genai.Image{
			ImageBytes: productBytes,
			MIMEType:   "image/jpeg",
		},
	}

	resp, err := client.Models.RecontextImage(ctx,
		modelName,
		&genai.RecontextImageSource{
			PersonImage:   personImage,
			ProductImages: []*genai.ProductImage{productImage},
		}, nil,
	)
	if err != nil {
		return fmt.Errorf("recontext image failed: %w", err)
	}

	// Ensure we have a generated image and save it
	if len(resp.GeneratedImages) == 0 || resp.GeneratedImages[0].Image == nil {
		return fmt.Errorf("no image was generated")
	}

	// TODO(developer): Update below lines
	path := "testdata"
	outputFile := "man_in_sweater.png"

	// Save output
	outputPath := filepath.Join(path, outputFile)

	if err := os.WriteFile(outputPath, resp.GeneratedImages[0].Image.ImageBytes, 0o644); err != nil {
		return fmt.Errorf("failed to save generated image: %w", err)
	}

	fmt.Fprintf(w, "Created output image using %d bytes\n", len(resp.GeneratedImages[0].Image.ImageBytes))
	fmt.Fprintln(w, outputPath)

	// Example response:
	// Created output image using 1636301 bytes ...

	return nil
}

// [END googlegenaisdk_imggen_virtual_try_on_with_txt_img]
