// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package snippets

// [START generativeaionvertexai_gemini_generate_from_text_input]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

func textInput(w io.Writer, projectID string, location string, modelName string) error {
	// location := "us-central1"
	// modelName := "gemini-1.0-pro-vision-001"

	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	gemini := client.GenerativeModel(modelName)
	// Does the returned sentiment score match the reviewer's movie rating?
	prompt := genai.Text(`Give a score from 1 - 10 to suggest if the
            following movie review is negative or positive (1 is most
            negative, 10 is most positive, 5 will be neutral). Include an
            explanation.

            The movie takes some time to build, but that is part of its beauty.
            By the time you are hooked, this tale of friendship and hope is
            thrilling and affecting, until the very last scene. You will find
            yourself rooting for the hero every step of the way. This is the
            sharpest, most original animated film I have seen in years.
            I would give it 8 out of 10 stars.`)

	resp, err := gemini.GenerateContent(ctx, prompt)
	if err != nil {
		return fmt.Errorf("error generating content: %w", err)
	}
	rb, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}
	fmt.Fprintln(w, string(rb))
	return nil
}

// [END generativeaionvertexai_gemini_generate_from_text_input]
