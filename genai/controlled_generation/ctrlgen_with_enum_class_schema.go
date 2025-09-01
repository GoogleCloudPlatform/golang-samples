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

// [START googlegenaisdk_ctrlgen_with_enum_class_schema]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

type InstrumentClass struct {
	Percussion string
	String     string
	Woodwind   string
	Brass      string
	Keyboard   string
}

var Instruments = InstrumentClass{
	Percussion: "Percussion",
	String:     "String",
	Woodwind:   "Woodwind",
	Brass:      "Brass",
	Keyboard:   "Keyboard",
}

// generateWithEnumClassSchema shows how to use a class-like enum schema to generate output.
func generateWithEnumClassSchema(w io.Writer) error {
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
			{Text: "What type of instrument is a guitar?"},
		}, Role: "user"},
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "text/x.enum",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeString,
			Enum: []string{
				Instruments.Percussion,
				Instruments.String,
				Instruments.Woodwind,
				Instruments.Brass,
				Instruments.Keyboard,
			},
		},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()
	fmt.Fprintln(w, respText)

	// Example response:
	// String

	return nil
}

// [END googlegenaisdk_ctrlgen_with_enum_class_schema]
