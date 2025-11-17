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

// Package count_tokens shows examples of counting tokens using the GenAI SDK.
package count_tokens

// [START googlegenaisdk_counttoken_localtokenizer_with_txt]
import (
	"fmt"
	"io"

	"google.golang.org/genai"
	"google.golang.org/genai/tokenizer"
)

// countTokenLocalWithTxt shows how to count tokens using the local tokenizer with text input.
func countTokenLocalWithTxt(w io.Writer) error {
	modelName := "gemini-2.5-flash"
	client, err := tokenizer.NewLocalTokenizer(modelName)
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: "What's the highest mountain in Africa?"},
		}},
	}

	resp, err := client.CountTokens(contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	fmt.Fprintf(w, "Total: %d\n", resp.TotalTokens)

	// Example response:
	// Total: 9

	return nil
}

// [END googlegenaisdk_counttoken_localtokenizer_with_txt]
