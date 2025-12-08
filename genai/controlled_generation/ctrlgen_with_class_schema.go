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

// [START googlegenaisdk_ctrlgen_with_class_schema]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// Recipe represents the schema for a recipe response.
type RecipeClass struct {
	RecipeName  string   `json:"recipe_name"`
	Ingredients []string `json:"ingredients"`
}

// generateWithClassSchema shows how to use class schema to generate output.
func generateWithClassSchema(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash"
	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: "List a few popular cookie recipes."},
		}, Role: "user"},
	}
	// JSON Schema for []Recipe
	schema := &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"recipe_name": {
					Type:        genai.TypeString,
					Description: "Name of the recipe",
				},
				"ingredients": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeString,
					},
					Description: "List of ingredients for the recipe",
				},
			},
			Required: []string{"recipe_name", "ingredients"},
		},
	}
	resp, err := client.Models.GenerateContent(ctx, modelName, contents, &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   schema, // Expect a list of Recipe objects
	})
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	// Print raw JSON string
	fmt.Fprintln(w, resp.Text())

	// Parse JSON into Go structs
	var recipes []RecipeClass
	if err := json.Unmarshal([]byte(resp.Text()), &recipes); err != nil {
		return fmt.Errorf("failed to parse response JSON: %w", err)
	}

	// Print parsed objects
	for _, r := range recipes {
		fmt.Fprintf(w, "Recipe: %s\n", r.RecipeName)
		fmt.Fprintf(w, "Ingredients: %v\n", r.Ingredients)
	}

	// Example output:
	// [
	//   {
	//	  "recipe_name": "Chocolate Chip Cookies"
	//    "ingredients": [
	//   	"all-purpose flour",
	//  	"baking soda",
	// 		"salt",
	//		"unsalted butter",
	//		"granulated sugar",
	//		"brown sugar",
	//		"eggs",
	//		"vanilla extract",
	//		"chocolate chips"
	//    ],
	//  },
	//   ...
	// ]

	return nil
}

// [END googlegenaisdk_ctrlgen_with_class_schema]
