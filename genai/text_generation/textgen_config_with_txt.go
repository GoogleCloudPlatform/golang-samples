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

// [START googlegenaisdk_textgen_config_with_txt]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithConfig shows how to generate text using a custom configuration.
func generateWithConfig(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{})
	if err != nil {
		return fmt.Errorf("unable to create genai client: %w", err)
	}

	modelName := "gemini-2.0-flash-001"
	contents := genai.Text("Why is the sky blue?")
	config := &genai.GenerateContentConfig{
		Temperature:      genai.Ptr(0.0),
		CandidateCount:   genai.Ptr(int64(1)),
		ResponseMIMEType: "application/json",
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("unable to generate content: %w", err)
	}

	respText, err := resp.Text()
	if err != nil {
		return fmt.Errorf("unable to convert model response to text: %w", err)
	}
	fmt.Fprintln(w, respText)
	// Example response:
	// {
  //   "explanation": "The sky is blue due to a phenomenon called Rayleigh scattering ...
	// }

	return nil
}

// [END googlegenaisdk_textgen_config_with_txt]
