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

// [START googlegenaisdk_tools_google_search_with_txt]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithGoogleSearch shows how to generate text using Google Search.
func generateWithGoogleSearch(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash"
	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: "When is the next total solar eclipse in the United States?"},
		},
			Role: "user"},
	}
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{GoogleSearch: &genai.GoogleSearch{ExcludeDomains: []string{"example.com", "example.org"}}},
		},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)

	// Example response:
	// The next total solar eclipse in the United States will occur on March 30, 2033, but it will only ...

	return nil
}

// [END googlegenaisdk_tools_google_search_with_txt]
