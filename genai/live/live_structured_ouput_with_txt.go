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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
)

// CalendarEvent represents the structured output we want the model to produce.
type CalendarEvent struct {
	Name         string   `json:"name"`
	Date         string   `json:"date"`
	Participants []string `json:"participants"`
}

// generateStructuredOutputWithTxt demonstrates calling the model via the Vertex AI
// OpenAI-compatible Chat Completions endpoint and parsing structured JSON into a Go struct.
func generateStructuredOutputWithTxt(w io.Writer) error {
	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		return fmt.Errorf("environment variable GOOGLE_CLOUD_PROJECT must be set")
	}

	location := "global"
	// Use "openapi" to call Gemini via the OpenAI-compatible endpoint.
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

	baseURL := fmt.Sprintf("https://aiplatform.googleapis.com/v1/projects/%s/locations/%s/endpoints/%s",
		projectID, location, endpointID)
	if location != "global" {
		baseURL = fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/endpoints/%s",
			location, projectID, location, endpointID)
	}

	modelName := "google/gemini-3.8-flash"

	payload := map[string]any{
		"model": modelName,
		"messages": []map[string]string{
			{"role": "system", "content": "Extract the event information."},
			{"role": "user", "content": "Alice and Bob are going to a science fair on Friday."},
		},
		"response_format": map[string]any{
			"type": "json_schema",
			"json_schema": map[string]any{
				"name":   "CalendarEvent",
				"strict": true,
				"schema": map[string]any{
					"type":                 "object",
					"additionalProperties": false,
					"properties": map[string]any{
						"name":         map[string]any{"type": "string"},
						"date":         map[string]any{"type": "string"},
						"participants": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
					},
					"required": []string{"name", "date", "participants"},
				},
			},
		},
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("chat completions request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("chat completions failed (%s): %s", resp.Status, string(respBody))
	}

	var completion struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &completion); err != nil {
		return fmt.Errorf("failed to parse completion response: %w", err)
	}
	if len(completion.Choices) == 0 || completion.Choices[0].Message.Content == "" {
		return fmt.Errorf("empty response text")
	}

	var event CalendarEvent
	if err := json.Unmarshal([]byte(completion.Choices[0].Message.Content), &event); err != nil {
		return fmt.Errorf("Model output was not valid JSON. Raw output:\n%s\n", completion.Choices[0].Message.Content)
	}

	// Print parsed struct in the same friendly format.
	fmt.Fprintln(w, event)

	// Example expected output:
	// {science fair Friday [Alice Bob]}

	return nil
}

// [END googlegenaisdk_live_structured_ouput_with_txt]
