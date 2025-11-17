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

// [START googlegenaisdk_counttoken_localtokenizer_compute_with_txt]
import (
	"fmt"
	"io"

	"google.golang.org/genai"
	"google.golang.org/genai/tokenizer"
)

// countTokenLocalComputeWithTxt shows how to compute tokens using the local tokenizer with text input.
func countTokenLocalComputeWithTxt(w io.Writer) error {
	modelName := "gemini-2.5-flash"
	client, err := tokenizer.NewLocalTokenizer(modelName)
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: "What's the longest word in the English language?"},
		}},
	}

	resp, err := client.ComputeTokens(contents)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	for _, tokenInfo := range resp.TokensInfo {
		fmt.Fprintf(w, "Role: %s\n", tokenInfo.Role)
		fmt.Fprintf(w, "Token IDs: %v\n", tokenInfo.TokenIDs)
		fmt.Fprintf(w, "Tokens: %v\n", tokenInfo.Tokens)
	}

	// Example response:
	// Role: user
	// Token IDs: [3689 236789 236751 506 27801 3658 528 506 5422 5192 236881]
	// Tokens: [[87 104 97 116] [39] [115] ...

	return nil
}

// [END googlegenaisdk_counttoken_localtokenizer_compute_with_txt]
