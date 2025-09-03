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

// Package text_generation shows examples of generating text using the GenAI SDK.
package text_generation

// [START googlegenaisdk_textgen_async_with_txt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateWithTextAsyncStream shows how to stream a text generation response.
func generateWithTextAsyncStream(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash"

	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: "Compose a song about the adventures of a time-traveling squirrel."},
			},
		},
	}

	for resp, err := range client.Models.GenerateContentStream(ctx, modelName, contents, &genai.GenerateContentConfig{
		ResponseModalities: []string{"TEXT"}}) {
		if err != nil {
			return fmt.Errorf("failed to generate content: %w", err)
		}

		chunk := resp.Text()

		fmt.Fprintln(w, chunk)
	}

	// Example output (streamed piece by piece):
	// (Verse 1)
	// Pip was a squirrel, a regular chap,
	//Burying acorns, enjoying a nap.
	//One sunny morning, beneath the old pine,
	//He dug up a thing, incredibly fine.
	//A tiny contraption, with gears and a gleam,
	//It pulsed with a power, a
	// time-traveling dream.
	//He nudged it with curiosity, twitching his nose,
	//And *poof!* went the world, as everyone knows...
	//
	//(Chorus)
	//Oh, Pip the squirrel, with his bushy brown tail,
	//Through the time stream he'd often sail!
	// ...

	return nil
}

// [END googlegenaisdk_textgen_async_with_txt]
