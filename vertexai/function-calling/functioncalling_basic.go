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

// Function calling lets developers create a description of a function in their code, then pass
// that description to a language model in a request.
//
// This function call example involves 3 messages:
// - ask the model to generate a function call request
// - call the Open API service (simulated in this example)
// - ask the model to interpret the function call response
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

// functionCallsBasic opens a chat session and sends 2 messages to the model:
// - first, to convert a text into a structured function call request
// - second, to convert a structured function call response into natural language.
// Writes output of second call to w.
func functionCallsBasic(w io.Writer, projectID, location, modelName string) error {
	// location := "us-central1"
	// modelName := "gemini-1.5-flash-001"
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	// Build an OpenAPI schema, in memory
	params := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"location": {
				Type:        genai.TypeString,
				Description: "location",
			},
		},
	}
	fundecl := &genai.FunctionDeclaration{
		Name:        "getCurrentWeather",
		Description: "Get the current weather in a given location",
		Parameters:  params,
	}
	model.Tools = []*genai.Tool{
		{FunctionDeclarations: []*genai.FunctionDeclaration{fundecl}},
	}

	chat := model.StartChat()

	resp, err := chat.SendMessage(ctx, genai.Text("What's the weather like in Boston?"))
	if err != nil {
		return err
	}
	if len(resp.Candidates) == 0 ||
		len(resp.Candidates[0].Content.Parts) == 0 {
		return errors.New("empty response from model")
	}

	// The model has returned a function call to the declared function `getCurrentWeather`
	// with a value for the argument `location`.
	_, err = json.MarshalIndent(resp.Candidates[0].Content.Parts[0], "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}

	// In this example, we'll use synthetic data to simulate a response payload from an external API
	weather := map[string]string{
		"location":    "Boston",
		"temperature": "38",
		"description": "Partly Cloudy",
		"icon":        "partly-cloudy",
		"humidity":    "65",
		"wind":        "{\"speed\": \"10\", \"direction\": \"NW\"}",
	}
	weather_json, _ := json.Marshal(weather)

	// Create a function call response, to simulate the result of a call to a
	// real service
	funresp := &genai.FunctionResponse{
		Name: "getCurrentWeather",
		Response: map[string]any{
			"currentWeather": weather_json,
		},
	}
	_, err = json.MarshalIndent(funresp, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}

	// And provide the function call response to the model
	resp, err = chat.SendMessage(ctx, funresp)
	if err != nil {
		return err
	}
	if len(resp.Candidates) == 0 ||
		len(resp.Candidates[0].Content.Parts) == 0 {
		return errors.New("empty response from model")
	}

	// The model has taken the function call response as input, and has
	// reformulated the response to the user.
	content, err := json.MarshalIndent(resp.Candidates[0].Content.Parts[0], "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}

	fmt.Fprintf(w, "generated summary:\n%s\n", content)
	return nil
}

// [END generativeaionvertexai_gemini_function_calling]
