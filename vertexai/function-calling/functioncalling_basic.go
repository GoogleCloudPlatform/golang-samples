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

// Function calling allows the model to improve its responses with the use of external
// data sources.
package functioncalling

// [START generativeaionvertexai_gemini_function_calling]
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

// functionCalling demonstrates how to submit a prompt and a function declaration to the model,
// allowing it to suggest a call to the function to fetch external data. Returning this data
// to the model enables it to generate a text response that incorporates the data.
func functionCalling(w io.Writer, projectID, location, modelName string) error {
	// location = "us-central1"
	// modelName = "gemini-1.5-flash-002"
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("failed to create GenAI client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)
	// Set temperature to 0.0 for maximum determinism in function calling.
	model.SetTemperature(0.0)

	funcName := "getCurrentWeather"
	funcDecl := &genai.FunctionDeclaration{
		Name:        funcName,
		Description: "Get the current weather in a given location",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"location": {
					Type:        genai.TypeString,
					Description: "location",
				},
			},
			Required: []string{"location"},
		},
	}
	// Add the weather function to our model toolbox.
	model.Tools = []*genai.Tool{
		{
			FunctionDeclarations: []*genai.FunctionDeclaration{funcDecl},
		},
	}

	prompt := genai.Text("What's the weather like in Boston?")
	resp, err := model.GenerateContent(ctx, prompt)

	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}
	if len(resp.Candidates) == 0 {
		return errors.New("got empty response from model")
	} else if len(resp.Candidates[0].FunctionCalls()) == 0 {
		return errors.New("got no function call suggestions from model")
	}

	// In a production environment, consider adding validations for function names and arguments.
	for _, fnCall := range resp.Candidates[0].FunctionCalls() {
		fmt.Fprintf(w, "The model suggests to call the function %q with args: %v\n", fnCall.Name, fnCall.Args)
		// Example response:
		// The model suggests to call the function "getCurrentWeather" with args: map[location:Boston]
	}
	// Use synthetic data to simulate a response from the external API.
	// In a real application, this would come from an actual weather API.
	mockAPIResp, err := json.Marshal(map[string]string{
		"location":         "Boston",
		"temperature":      "38",
		"temperature_unit": "F",
		"description":      "Cold and cloudy",
		"humidity":         "65",
		"wind":             `{"speed": "10", "direction": "NW"}`,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal function response to JSON: %w", err)
	}

	funcResp := &genai.FunctionResponse{
		Name: funcName,
		Response: map[string]any{
			"content": mockAPIResp,
		},
	}

	// Return the API response to the model allowing it to complete its response.
	resp, err = model.GenerateContent(ctx, prompt, funcResp)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return errors.New("got empty response from model")
	}

	fmt.Fprintln(w, resp.Candidates[0].Content.Parts[0])
	// Example response:
	// The weather in Boston is cold and cloudy, with a humidity of 65% and a temperature of 38°F. ...

	return nil
}

// [END generativeaionvertexai_gemini_function_calling]
