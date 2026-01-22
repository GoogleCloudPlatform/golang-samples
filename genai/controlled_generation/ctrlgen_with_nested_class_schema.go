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

// [START googlegenaisdk_ctrlgen_with_nested_class_schema]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

type Grade string

const (
	APlus Grade = "a+"
	A     Grade = "a"
	B     Grade = "b"
	C     Grade = "c"
	D     Grade = "d"
	F     Grade = "f"
)

type Recipe struct {
	RecipeName string `json:"recipe_name"`
	Rating     Grade  `json:"rating"`
}

// generateWithNestedClassSchema shows how to use nested class schema to generate output.
func generateWithNestedClassSchema(w io.Writer) error {
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
			{Text: "List about 10 home-baked cookies and give them grades based on tastiness."},
		}, Role: genai.RoleUser},
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"recipe_name": {Type: genai.TypeString},
					"rating": {
						Type: genai.TypeString,
						Enum: []string{string(APlus), string(A), string(B), string(C), string(D), string(F)},
					},
				},
				Required: []string{"recipe_name", "rating"},
			},
		},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	fmt.Fprintln(w, resp.Text())

	// Example response:
	// [
	//    {"rating":"a+","recipe_name":"Chocolate Chip Cookies"},
	// 	  {"rating":"a","recipe_name":"Oatmeal Raisin Cookies"},
	//	  {"rating":"a+","recipe_name":"Peanut Butter Cookies"},
	//	  {"rating":"b","recipe_name":"Snickerdoodle Cookies"},
	//   ...
	// ]

	return nil
}

// [END googlegenaisdk_ctrlgen_with_nested_class_schema]
