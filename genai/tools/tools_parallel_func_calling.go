// Copyright 2026 Google LLC
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

// [START aiplatform_genai_parallel_func_calling]
import (
	"context"
	"errors"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithParallelFunctionCalling shows how to execute multiple function calls in parallel
// and return their results to the model for generating a complete response.
func generateWithParallelFunctionCalling(w io.Writer) error {
	ctx := context.Background()

	// Initialize the unified GenAI client for Vertex AI.
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
		Backend:     genai.BackendVertexAI,
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	funcName := "getCurrentWeather"
	funcDecl := &genai.FunctionDeclaration{
		Name:        funcName,
		Description: "Get the current weather in a given location",
		Parameters: &genai.Schema{
			Type: "object",
			Properties: map[string]*genai.Schema{
				"location": {
					Type: "string",
					Description: "The location for which to get the weather. " +
						"It can be a city name, a city name and state, or a zip code. " +
						"Examples: 'San Francisco', 'San Francisco, CA', '95616', etc.",
				},
			},
			Required: []string{"location"},
		},
	}

	config := &genai.GenerateContentConfig{
		Temperature: genai.Ptr(float32(0.0)),
		Tools: []*genai.Tool{
			{
				FunctionDeclarations: []*genai.FunctionDeclaration{funcDecl},
			},
		},
	}

	// Initialize the conversation history with the user prompt.
	prompt := "Get weather details in New Delhi and San Francisco?"
	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: prompt},
			},
		},
	}

	modelName := "gemini-2.5-flash"

	// First API call: The model determines it needs to call tools based on the prompt.
	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return errors.New("got empty response from model")
	}

	// Extract the parallel function call requests from the model's response.
	var functionCalls []*genai.FunctionCall
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.FunctionCall != nil {
			functionCalls = append(functionCalls, part.FunctionCall)
			fmt.Fprintf(w, "Model suggests to call the function %q with args: %v\n", part.FunctionCall.Name, part.FunctionCall.Args)
		}
	}

	if len(functionCalls) == 0 {
		return errors.New("got no function call suggestions from model")
	}

	// Append the model's tool call request to the conversation history.
	contents = append(contents, resp.Candidates[0].Content)

	// Simulate external API responses. The SDK now directly accepts map[string]any.
	mockAPIResp1 := map[string]any{
		"location":         "New Delhi",
		"temperature":      "42",
		"temperature_unit": "C",
		"description":      "Hot and humid",
		"humidity":         "65",
	}

	mockAPIResp2 := map[string]any{
		"location":         "San Francisco",
		"temperature":      "36",
		"temperature_unit": "F",
		"description":      "Cold and cloudy",
		"humidity":         "N/A",
	}

	// Bundle the API responses into a single Content block as Parts.
	funcRespContent := &genai.Content{
		Role: "user",
		Parts: []*genai.Part{
			{
				FunctionResponse: &genai.FunctionResponse{
					Name:     funcName,
					Response: mockAPIResp1,
				},
			},
			{
				FunctionResponse: &genai.FunctionResponse{
					Name:     funcName,
					Response: mockAPIResp2,
				},
			},
		},
	}

	// Append the tool response to the conversation history
	contents = append(contents, funcRespContent)

	// Final API call: The model synthesizes the tool results into a natural language response.
	resp, err = client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return errors.New("got empty response from model")
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			fmt.Fprintln(w, part.Text)
		}
	}

	// Example response:
	// The weather in New Delhi is hot and humid with a temperature of 42 degrees Celsius. The weather in San Francisco is ...

	return nil
}

// [END aiplatform_genai_parallel_func_calling]
