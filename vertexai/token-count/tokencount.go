// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// token-count shows an example of determining how many tokens correspond to
// a given prompt string
package tokencount

// [START aiplatform_gemini_token_count]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

// countTokens returns the number of tokens for this prompt.
func countTokens(w io.Writer, prompt, projectID, location, modelName string) error {
	// prompt := "why is the sky blue?"
	// location := "us-central1"
	// modelName := "gemini-1.0-pro-002"

	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("unable to create client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	resp, err := model.CountTokens(ctx, genai.Text(prompt))
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Number of tokens for the prompt: %d\n", resp.TotalTokens)

	return nil
}

// [END aiplatform_gemini_token_count]
