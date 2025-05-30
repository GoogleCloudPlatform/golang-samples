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

// [START googlegenaisdk_counttoken_resp_with_txt]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateTextAndCount shows how to generate text and obtain token count metadata from the model response.
func generateTextAndCount(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.0-flash-001"
	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: "Why is the sky blue?"},
		},
			Role: "user"},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	usage, err := json.MarshalIndent(resp.UsageMetadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to convert usage metadata to JSON: %w", err)
	}
	fmt.Fprintln(w, string(usage))

	// Example response:
	// {
	// 	 "candidatesTokenCount": 339,
	// 	 "promptTokenCount": 6,
	// 	 "totalTokenCount": 345
	// }

	return nil
}

// [END googlegenaisdk_counttoken_resp_with_txt]
