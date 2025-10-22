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

// [START googlegenaisdk_tools_urlcontext_with_txt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateURLContentWithText demonstrates using the URL Context tool
// to compare and reason over the content of external URLs.
func generateURLContentWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// Model and input prompt
	modelName := "gemini-2.5-flash"
	// TODO(developer): Replace with your own URLs
	url1 := "https://cloud.google.com/vertex-ai/docs/generative-ai/start"
	url2 := "https://cloud.google.com/docs/overview"

	prompt := fmt.Sprintf("Compare the content, purpose, and audiences of %s and %s.", url1, url2)

	// Build the request configuration with the URL Context tool
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{
				URLContext: &genai.URLContext{},
			},
		},
		ResponseModalities: []string{"TEXT"},
	}

	// Generate content using the model
	resp, err := client.Models.GenerateContent(ctx, modelName, []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: prompt},
			},
		},
	}, config)
	if err != nil {
		return fmt.Errorf("generate content failed: %w", err)
	}

	// Print the model output
	fmt.Fprintln(w, resp.Text())

	// Optionally, print retrieved URL metadata
	if len(resp.Candidates) > 0 && resp.Candidates[0].URLContextMetadata != nil {
		fmt.Fprintf(w, "\nRetrieved URL metadata: %+v\n", resp.Candidates[0].URLContextMetadata)
	}

	// Example output:
	// Here's an analysis of "https://cloud.google.com/docs/overview":
	//
	//*   **Content:** This page provides a high-level overview of Google Cloud
	//...

	return nil
}

// [END googlegenaisdk_tools_urlcontext_with_txt]
