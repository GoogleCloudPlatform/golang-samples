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

// [START googlegenaisdk_textgen_with_txt_stream]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithTextStream shows how to generate text stream using a text prompt.
func generateWithTextStream(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash"
	contents := genai.Text("Why is the sky blue?")

	for resp, err := range client.Models.GenerateContentStream(ctx, modelName, contents, nil) {
		if err != nil {
			return fmt.Errorf("failed to generate content: %w", err)
		}

		chunk := resp.Text()

		fmt.Fprintln(w, chunk)
	}

	// Example response:
	// The
	//  sky is blue
	//  because of a phenomenon called **Rayleigh scattering**. Here's the breakdown:
	// ...

	return nil
}

// [END googlegenaisdk_textgen_with_txt_stream]
