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

// Package express_mode shows how to use the GenAI SDK to generate text with VertexAI Api Key.
package express_mode

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

var newClient = genai.NewClient

// generateContent shows how to use Vertex AI Express mode with an API key.
func generateContentWithApiKey(w io.Writer) error {
	ctx := context.Background()

	// TODO(developer): Replace with your actual API key
	apiKey := "YOUR_API_KEY"

	client, err := newClient(ctx, &genai.ClientConfig{
		APIKey:      apiKey,
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash"
	resp, err := client.Models.GenerateContent(ctx, modelName,
		[]*genai.Content{
			{Parts: []*genai.Part{
				{Text: "Explain bubble sort to me."},
			}, Role: "user"},
		},
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	fmt.Fprintln(w, resp.Text())

	// Example response:
	// Bubble Sort is a simple sorting algorithm that repeatedly steps through the list

	return nil
}

// [END googlegenaisdk_vertexai_express_mode]
