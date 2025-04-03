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

// [START generativeaionvertexai_gemini_controlled_generation_response_schema_6]
import (
	"context"
	"errors"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

// controlledGenerationResponseSchema6 shows how to make sure the generated output
// will always be valid JSON and adhere to a specific schema.
func controlledGenerationResponseSchema6(w io.Writer, projectID, location, modelName string) error {
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
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"object": {
						Type: genai.TypeString,
					},
				},
			},
		},
	}

	// These images in Cloud Storage are viewable at
	// https://storage.googleapis.com/cloud-samples-data/generative-ai/image/office-desk.jpeg
	// https://storage.googleapis.com/cloud-samples-data/generative-ai/image/gardening-tools.jpeg

	img1 := genai.FileData{
		MIMEType: "image/jpeg",
		FileURI:  "gs://cloud-samples-data/generative-ai/image/office-desk.jpeg",
	}

	img2 := genai.FileData{
		MIMEType: "image/jpeg",
		FileURI:  "gs://cloud-samples-data/generative-ai/image/gardening-tools.jpeg",
	}

	prompt := "Generate a list of objects in the images."

	res, err := model.GenerateContent(ctx, img1, img2, genai.Text(prompt))
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

// [END generativeaionvertexai_gemini_controlled_generation_response_schema_6]
