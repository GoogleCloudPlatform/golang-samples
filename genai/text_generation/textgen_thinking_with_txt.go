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

// [START googlegenaisdk_thinking_textgen_with_txt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateThinkingWithText shows how to generate thinking using a text prompt.
func generateThinkingWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	resp, err := client.Models.GenerateContent(ctx,
		"gemini-2.5-flash",
		genai.Text("solve x^2 + 4x + 4 = 0"),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)
	// Example response:
	// To solve the quadratic equation $x^2 + 4x + 4 = 0$, we can use a few methods:
	//
	// **Method 1: Factoring (Recognizing a Perfect Square Trinomial)**
	// **1. The Foundation: Data and Algorithms**
	//
	// Notice that the left side of the equation is a perfect square trinomial.
	// ...

	return nil
}

// [END googlegenaisdk_thinking_textgen_with_txt]
