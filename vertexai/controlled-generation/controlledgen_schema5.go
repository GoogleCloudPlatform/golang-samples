// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controlledgeneration

// [START generativeaionvertexai_gemini_controlled_generation_response_schema_5]
import (
	"context"
	"errors"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

// controlledGenerationResponseSchema5 shows how to make sure the generated output
// will always be valid JSON and adhere to a specific schema.
func controlledGenerationResponseSchema5(w io.Writer, projectID, location, modelName string) error {
	// location := "us-central1"
	// modelName := "gemini-1.5-pro-001"
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	model.GenerationConfig.ResponseMIMEType = "application/json"

	// Build an OpenAPI schema, in memory
	model.GenerationConfig.ResponseSchema = &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"AnnouncementDate": {
					Type: genai.TypeString,
				},
				"Authors": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeString,
					},
				},
				"JournalRef": {
					Type: genai.TypeString,
				},
				"Keywords": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeString,
					},
				},
				"SubmissionDate": {
					Type: genai.TypeString,
				},
				"Title": {
					Type: genai.TypeString,
				},
				"Version": {
					Type: genai.TypeString,
					Enum: []string{
						"Dungeons & Dragons",
						"Duel Masters",
						"G.I. Joe",
						"Jem and The Holograms",
						"Littlest Pet Shop",
						"Magic: The Gathering",
						"Monopoly",
						"My Little Pony",
						"Nerf",
					},
				},
			},
		},
	}

	prompt := `
		Cymbal stock slid 5.2% following a double-downgrade to “underperform” from “buy” at Bank of Anthos.
		Bank of Anthos conducted a “deep dive” on trading card game. Bank of Anthos said Cymbal has been
		overprinting cards and destroying the long-term value of the business.
	`

	res, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return fmt.Errorf("unable to generate contents: %v", err)
	}

	if len(res.Candidates) == 0 ||
		len(res.Candidates[0].Content.Parts) == 0 {
		return errors.New("empty response from model")
	}

	fmt.Fprint(w, res.Candidates[0].Content.Parts[0])
	return nil
}

// [END generativeaionvertexai_gemini_controlled_generation_response_schema_5]
