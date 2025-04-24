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

// [START generativeaionvertexai_gemini_controlled_generation_response_schema_7]
import (
	"context"
	"errors"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

// controlledGenerationResponseSchemaEnum demonstrates how to constrain model responses
// to a predefined set of enum values for genre classification.
func controlledGenerationResponseSchemaEnum(w io.Writer, projectID, location, modelName string) error {
	// location = "us-central1"
	// modelName = "gemini-2.0-flash-001"
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("failed to create GenAI client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	model.GenerationConfig.ResponseMIMEType = "text/x.enum"
	model.GenerationConfig.ResponseSchema = &genai.Schema{
		Type: genai.TypeString,
		Enum: []string{"drama", "comedy", "documentary"},
	}

	prompt := `
The film aims to educate and inform viewers about real-life subjects, events, or people.
It offers a factual record of a particular topic by combining interviews, historical footage,
and narration. The primary purpose of a film is to present information and provide insights
into various aspects of reality.
`

	res, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return errors.New("got empty response from model")
	}

	fmt.Fprintf(w, "Candidate label: %q", res.Candidates[0].Content.Parts[0])
	// Example response:
	// Candidate label: "documentary"

	return nil
}

// [END generativeaionvertexai_gemini_controlled_generation_response_schema_7]
