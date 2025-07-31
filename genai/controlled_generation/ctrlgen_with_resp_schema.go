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

// Package controlled_generation shows how to use the GenAI SDK to generate text that adheres to a specific schema.
package controlled_generation

// [START googlegenaisdk_ctrlgen_with_resp_schema]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithRespSchema shows how to use a response schema to generate output in a specific format.
func generateWithRespSchema(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		// See the OpenAPI specification for more details and examples:
		//   https://spec.openapis.org/oas/v3.0.3.html#schema-object
		ResponseSchema: &genai.Schema{
			Type: "array",
			Items: &genai.Schema{
				Type: "object",
				Properties: map[string]*genai.Schema{
					"recipe_name": {Type: "string"},
					"ingredients": {
						Type:  "array",
						Items: &genai.Schema{Type: "string"},
					},
				},
				Required: []string{"recipe_name", "ingredients"},
			},
		},
	}
	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: "List a few popular cookie recipes."},
		},
			Role: "user"},
	}
	modelName := "gemini-2.5-flash"

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)

	// Example response:
	// [
	//   {
	//     "ingredients": [
	//       "2 1/4 cups all-purpose flour",
	//       "1 teaspoon baking soda",
	//       ...
	//     ],
	//     "recipe_name": "Chocolate Chip Cookies"
	//   },
	//   {
	//     ...
	//   },
	//   ...
	// ]

	return nil
}

// [END googlegenaisdk_ctrlgen_with_resp_schema]
