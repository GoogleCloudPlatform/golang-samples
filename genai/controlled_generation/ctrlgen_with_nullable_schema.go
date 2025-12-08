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

// [START googlegenaisdk_ctrlgen_with_nullable_schema]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithNullables shows how to use the response schema with nullable values.
func generateWithNullables(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash"
	prompt := `
The week ahead brings a mix of weather conditions.
Sunday is expected to be sunny with a temperature of 77°F and a humidity level of 50%. Winds will be light at around 10 km/h.
Monday will see partly cloudy skies with a slightly cooler temperature of 72°F and the winds will pick up slightly to around 15 km/h.
Tuesday brings rain showers, with temperatures dropping to 64°F and humidity rising to 70%.
Wednesday may see thunderstorms, with a temperature of 68°F.
Thursday will be cloudy with a temperature of 66°F and moderate humidity at 60%.
Friday returns to partly cloudy conditions, with a temperature of 73°F and the Winds will be light at 12 km/h.
Finally, Saturday rounds off the week with sunny skies, a temperature of 80°F, and a humidity level of 40%. Winds will be gentle at 8 km/h.
`
	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: prompt},
		},
			Role: "user"},
	}
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		// See the OpenAPI specification for more details and examples:
		//   https://spec.openapis.org/oas/v3.0.3.html#schema-object
		ResponseSchema: &genai.Schema{
			Type: "object",
			Properties: map[string]*genai.Schema{
				"forecast": {
					Type: "array",
					Items: &genai.Schema{
						Type: "object",
						Properties: map[string]*genai.Schema{
							"Day":         {Type: "string", Nullable: genai.Ptr(true)},
							"Forecast":    {Type: "string", Nullable: genai.Ptr(true)},
							"Temperature": {Type: "integer", Nullable: genai.Ptr(true)},
							"Humidity":    {Type: "string", Nullable: genai.Ptr(true)},
							"Wind Speed":  {Type: "integer", Nullable: genai.Ptr(true)},
						},
						Required: []string{"Day", "Temperature", "Forecast", "Wind Speed"},
					},
				},
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
	// {
	// 	"forecast": [
	// 		{"Day": "Sunday", "Forecast": "Sunny", "Temperature": 77, "Wind Speed": 10, "Humidity": "50%"},
	// 		{"Day": "Monday", "Forecast": "Partly Cloudy", "Temperature": 72, "Wind Speed": 15},
	// 		{"Day": "Tuesday", "Forecast": "Rain Showers", "Temperature": 64, "Wind Speed": null, "Humidity": "70%"},
	// 		{"Day": "Wednesday", "Forecast": "Thunderstorms", "Temperature": 68, "Wind Speed": null},
	// 		{"Day": "Thursday", "Forecast": "Cloudy", "Temperature": 66, "Wind Speed": null, "Humidity": "60%"},
	// 		{"Day": "Friday", "Forecast": "Partly Cloudy", "Temperature": 73, "Wind Speed": 12},
	// 		{"Day": "Saturday", "Forecast": "Sunny", "Temperature": 80, "Wind Speed": 8, "Humidity": "40%"}
	// 	]
	// }

	return nil
}

// [END googlegenaisdk_ctrlgen_with_nullable_schema]
