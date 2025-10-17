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

// [START googlegenaisdk_tools_google_search_and_urlcontext_with_txt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateGoogleSearchAndUrlContextWithTxt demonstrates using both
// the Google Search Tool and the URL Context Tool with Gemini.
func generateGSearchURLContentWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1beta1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash"
	// Define both tools: URL Context and Google Search
	tools := []*genai.Tool{
		{URLContext: &genai.URLContext{}},
		{GoogleSearch: &genai.GoogleSearch{}},
	}

	// TODO(developer): Replace with your own URL
	url := "https://www.google.com/search?q=events+in+New+York"

	prompt := fmt.Sprintf(
		"Give me three day events schedule based on %s. "+
			"Also let me know what needs to be taken care of considering weather and commute.",
		url,
	)

	// Build the generation config
	config := &genai.GenerateContentConfig{
		Tools:              tools,
		ResponseModalities: []string{"TEXT"},
	}

	// Call the model
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

	// Print the response text
	fmt.Fprintln(w, resp.Text())

	// Print retrieved URLs metadata if available
	if len(resp.Candidates) > 0 && resp.Candidates[0].URLContextMetadata != nil {
		fmt.Fprintf(w, "\nRetrieved URL metadata: %+v\n", resp.Candidates[0].URLContextMetadata)
	}

	// Example output:
	// Here is a three-day event schedule for New York City from Friday, October 17, 2025, to Sunday, October 19, 2025, along with weather and commute considerations.
	//
	//**Weather Forecast (October 17-19, 2025):**
	//*   **Friday, October 17:** Sunny and slightly breezy, with highs in the mid to upper 50s°F (13-15°C) and lows in the mid 30s°F (1-4°C). No rain is expected.
	//*   **Saturday, October 18:** A mix of sun and clouds, possibly clear. Highs in the lower 60s°F (around 16°C) and lows around 41-57°F (5-14°C). No rain is expected.
	//*   **Sunday, October 19:** Clear with highs around 66°F (19°C) and lows around 44-46°F (7-8°C). There's a slight chance of very light rain (0.01 inch).
	//...

	return nil
}

// [END googlegenaisdk_tools_google_search_and_urlcontext_with_txt]
