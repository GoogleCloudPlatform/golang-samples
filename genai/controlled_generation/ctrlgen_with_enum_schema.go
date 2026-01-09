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

// [START googlegenaisdk_ctrlgen_with_enum_schema]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithEnumSchema shows how to use enum schema to generate output.
func generateWithEnumSchema(w io.Writer) error {
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
			{Text: "What type of instrument is an oboe?"},
		}, Role: genai.RoleUser},
	}
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "text/x.enum",
		ResponseSchema: &genai.Schema{
			Type: "STRING",
			Enum: []string{"Percussion", "String", "Woodwind", "Brass", "Keyboard"},
		},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)

	// Example response:
	// Woodwind

	return nil
}

// [END googlegenaisdk_ctrlgen_with_enum_schema]
