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

// Package live shows how to use the GenAI SDK to generate text with live resources.
package live

// [START googlegenaisdk_live_structured_ouput_with_txt]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/genai"
)

// CalendarEvent represents the structured output we want the model to produce.
type CalendarEvent struct {
	Name         string   `json:"name"`
	Date         string   `json:"date"`
	Participants []string `json:"participants"`
}

// generateStructuredOutputWithTxt demonstrates calling the model via an OpenAPI-style
// endpoint (base_url + api_key) and parsing JSON output into a Go struct.
func generateStructuredOutputWithTxt(w io.Writer) error {
	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		return fmt.Errorf("environment variable GOOGLE_CLOUD_PROJECT must be set")
	}

	location := "us-central1"
	// Use "openapi" to call the Gemini API via the OpenAPI-compatible endpoint.
	endpointID := "openapi"

	// Programmatically obtain an access token.
	ts, err := google.DefaultTokenSource(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return fmt.Errorf("failed to get default token source: %w", err)
	}
	token, err := ts.Token()
	if err != nil {
		return fmt.Errorf("failed to fetch token: %w", err)
	}
	apiKey := token.AccessToken

	baseURL := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/endpoints/%s",
		location, projectID, location, endpointID)

	// Create a genai client that points to the OpenAPI endpoint, authenticating with the token.
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
		HTTPOptions: genai.HTTPOptions{
			BaseURL:    baseURL,
			APIVersion: "v1",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// Build the messages (system + user)
	contents := []*genai.Content{
		{
			Role: "system",
			Parts: []*genai.Part{
				{Text: "Extract the event information."},
			},
		},
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: "Alice and Bob are going to a science fair on Friday."},
			},
		},
	}

	// Ask the model to return JSON by setting a strict instruction and also request JSON mime type
	// to encourage machine-readable output.
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
	}

	modelName := "google/gemini-2.0-flash-001"

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("generate content failed: %w", err)
	}

	// Resp.Text() returns concatenated text of the top candidate.
	respText := resp.Text()
	if respText == "" {
		return fmt.Errorf("empty response text")
	}

	// Try to parse the JSON into our struct.
	var event CalendarEvent
	if err := json.Unmarshal([]byte(respText), &event); err != nil {
		fmt.Fprintf(w, "Model output was not valid JSON. Raw output:\n%s\n", respText)
		return nil
	}

	// Print parsed struct in the same friendly format.
	fmt.Fprintln(w, event)

	// Example expected output:
	// Parsed struct: {Name:science fair Date:Friday Participants:[Alice Bob]}

	return nil
}

// [END googlegenaisdk_live_structured_ouput_with_txt]
