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

// Package tools shows examples of various tools that Gemini model can use to generate text.
package tools

// [START googlegenaisdk_tools_func_desc_with_txt]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithFuncCall shows how to submit a prompt and a function declaration to the model,
// allowing it to suggest a call to the function to fetch external data. Returning this data
// enables the model to generate a text response that incorporates the data.
func generateWithFuncCall(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	weatherFunc := &genai.FunctionDeclaration{
		Description: "Returns the current weather in a location.",
		Name:        "getCurrentWeather",
		Parameters: &genai.Schema{
			Type: "object",
			Properties: map[string]*genai.Schema{
				"location": {Type: "string"},
			},
			Required: []string{"location"},
		},
	}
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{FunctionDeclarations: []*genai.FunctionDeclaration{weatherFunc}},
		},
		Temperature: genai.Ptr(float32(0.0)),
	}

	modelName := "gemini-2.0-flash-001"
	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: "What is the weather like in Boston?"},
		},
			Role: "user"},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	var funcCall *genai.FunctionCall
	for _, p := range resp.Candidates[0].Content.Parts {
		if p.FunctionCall != nil {
			funcCall = p.FunctionCall
			fmt.Fprint(w, "The model suggests to call the function ")
			fmt.Fprintf(w, "%q with args: %v\n", funcCall.Name, funcCall.Args)
			// Example response:
			// The model suggests to call the function "getCurrentWeather" with args: map[location:Boston]
		}
	}
	if funcCall == nil {
		return fmt.Errorf("model did not suggest a function call")
	}

	// Use synthetic data to simulate a response from the external API.
	// In a real application, this would come from an actual weather API.
	funcResp := &genai.FunctionResponse{
		Name: "getCurrentWeather",
		Response: map[string]any{
			"location":         "Boston",
			"temperature":      "38",
			"temperature_unit": "F",
			"description":      "Cold and cloudy",
			"humidity":         "65",
			"wind":             `{"speed": "10", "direction": "NW"}`,
		},
	}

	// Return conversation turns and API response to complete the model's response.
	contents = []*genai.Content{
		{Parts: []*genai.Part{
			{Text: "What is the weather like in Boston?"},
		},
			Role: "user"},
		{Parts: []*genai.Part{
			{FunctionCall: funcCall},
		}},
		{Parts: []*genai.Part{
			{FunctionResponse: funcResp},
		}},
	}

	resp, err = client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)

	// Example response:
	// The weather in Boston is cold and cloudy with a temperature of 38 degrees Fahrenheit. The humidity is ...

	return nil
}

// [END googlegenaisdk_tools_func_desc_with_txt]
