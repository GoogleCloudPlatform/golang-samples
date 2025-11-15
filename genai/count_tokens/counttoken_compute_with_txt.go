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

// [START googlegenaisdk_counttoken_compute_with_txt]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// computeWithTxt shows how to compute tokens with text input.
func computeWithTxt(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash"
	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: "What's the longest word in the English language?"},
		},
			Role: genai.RoleUser},
	}

	resp, err := client.Models.ComputeTokens(ctx, modelName, contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	type tokenInfoDisplay struct {
		IDs    []int64  `json:"token_ids"`
		Tokens []string `json:"tokens"`
	}
	// See the documentation: https://pkg.go.dev/google.golang.org/genai#ComputeTokensResponse
	for _, instance := range resp.TokensInfo {
		display := tokenInfoDisplay{
			IDs:    instance.TokenIDs,
			Tokens: make([]string, len(instance.Tokens)),
		}
		for i, t := range instance.Tokens {
			display.Tokens[i] = string(t)
		}

		data, err := json.MarshalIndent(display, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal token info: %w", err)
		}
		fmt.Fprintln(w, string(data))
	}

	// Example response:
	// {
	// 	"ids": [
	// 		1841,
	// 		235303,
	// 		235256,
	//    ...
	// 	],
	// 	"values": [
	// 		"What",
	// 		"'",
	// 		"s",
	//    ...
	// 	]
	// }

	return nil
}

// [END googlegenaisdk_counttoken_compute_with_txt]
