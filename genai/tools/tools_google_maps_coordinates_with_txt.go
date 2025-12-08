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

// [START googlegenaisdk_tools_google_maps_coordinates_with_txt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateGoogleMapsCoordinatesWithText demonstrates using the Google Maps tool
// to get responses localized to specific coordinates.
func generateGoogleMapsCoordinatesWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// Model and input prompt
	modelName := "gemini-2.5-flash"
	prompt := "Where can I get the best espresso near me?"
	lat := 40.7128
	long := -74.0060

	// Build the request configuration with the Google Maps tool
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{
				GoogleMaps: &genai.GoogleMaps{},
			},
		},
		ToolConfig: &genai.ToolConfig{
			RetrievalConfig: &genai.RetrievalConfig{
				LatLng: &genai.LatLng{
					Latitude:  &lat,  // New York City
					Longitude: &long, // NYC longitude
				},
				LanguageCode: "en_US", // Optional: localize Maps results
			},
		},
	}

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

	fmt.Fprintln(w, resp.Text())

	// Example output:
	// Here are some of the top-rated coffee shops near you that serve espresso:..."
	return nil
}

// [END googlegenaisdk_tools_google_maps_coordinates_with_txt]
