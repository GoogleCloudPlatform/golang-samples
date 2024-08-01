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

// [START generativeaionvertexai_gemini_controlled_generation_response_schema_3]
import (
	"context"
	"errors"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

// controlledGenerationResponseSchema3 shows how to make sure the generated output
// will always be valid JSON and adhere to a specific schema.
func controlledGenerationResponseSchema3(w io.Writer, projectID, location, modelName string) error {
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
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"forecast": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"Day": {
							Type: genai.TypeString,
						},
						"Forecast": {
							Type: genai.TypeString,
						},
						"Humidity": {
							Type: genai.TypeString,
						},
						"Temperature": {
							Type: genai.TypeInteger,
						},
						"Wind Speed": {
							Type: genai.TypeInteger,
						},
					},
					Required: []string{"Day", "Temperature", "Forecast"},
				},
			},
		},
	}

	prompt := `
		The week ahead brings a mix of weather conditions.
		Sunday is expected to be sunny with a temperature of 77°F and a humidity level of 50%. Winds will be light at around 10 km/h.
		Monday will see partly cloudy skies with a slightly cooler temperature of 72°F and humidity increasing to 55%. Winds will pick up slightly to around 15 km/h.
		Tuesday brings rain showers, with temperatures dropping to 64°F and humidity rising to 70%. Expect stronger winds at 20 km/h.
		Wednesday may see thunderstorms, with a temperature of 68°F and high humidity of 75%. Winds will be gusty at 25 km/h.
		Thursday will be cloudy with a temperature of 66°F and moderate humidity at 60%. Winds will ease slightly to 18 km/h.
		Friday returns to partly cloudy conditions, with a temperature of 73°F and lower humidity at 45%. Winds will be light at 12 km/h.
		Finally, Saturday rounds off the week with sunny skies, a temperature of 80°F, and a humidity level of 40%. Winds will be gentle at 8 km/h.
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

// [END generativeaionvertexai_gemini_controlled_generation_response_schema_3]
