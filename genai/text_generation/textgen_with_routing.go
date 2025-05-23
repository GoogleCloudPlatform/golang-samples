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

// [START googlegenaisdk_textgen_with_routing]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateWithRouting shows how to generate text using a text prompt and routing configuration.
func generateWithRouting(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1beta1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	generateContentConfig := &genai.GenerateContentConfig{ModelSelectionConfig: &genai.ModelSelectionConfig{FeatureSelectionPreference: genai.FeatureSelectionPreferencePrioritizeQuality}}

	resp, err := client.Models.GenerateContent(ctx,
		"model-optimizer-exp-04-09",
		genai.Text("How does AI work?"),
		generateContentConfig,
	)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()
	if err != nil {
		return fmt.Errorf("failed to convert model response to text: %w", err)
	}
	fmt.Fprintln(w, respText)
	// Example response:
	// That's a great question! Understanding how AI works can feel like ...
	// ...
	// **1. The Foundation: Data and Algorithms**
	// ...

	return nil
}

// [END googlegenaisdk_textgen_with_routing]
