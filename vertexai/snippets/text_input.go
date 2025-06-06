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

func generateContentFromText(w io.Writer, projectID string) error {
	location := "us-central1"
	modelName := "gemini-2.0-flash-001"

	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	gemini := client.GenerativeModel(modelName)
	prompt := genai.Text(
		"What's a good name for a flower shop that specializes in selling bouquets of dried flowers?")

	resp, err := gemini.GenerateContent(ctx, prompt)
	if err != nil {
		return fmt.Errorf("error generating content: %w", err)
	}
	// See the JSON response in
	// https://pkg.go.dev/cloud.google.com/go/vertexai/genai#GenerateContentResponse.
	rb, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}
	fmt.Fprintln(w, string(rb))
	return nil
}

// [END generativeaionvertexai_gemini_generate_from_text_input]
