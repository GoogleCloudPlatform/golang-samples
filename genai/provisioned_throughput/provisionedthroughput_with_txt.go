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

// Package provisionedthroughput shows examples of Gemini model can use to generate with text.
package provisionedthroughput

// [START googlegenaisdk_provisionedthroughput_with_txt]
import (
	"context"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/genai"
)

// generateProvisionedThroughputWithText shows how to generate text Provisioned Throughput.
func generateProvisionedThroughputWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{
			APIVersion: "v1",
			Headers: http.Header{
				// Options:
				// - "dedicated": Use Provisioned Throughput
				// - "shared": Use pay-as-you-go
				// https://cloud.google.com/vertex-ai/generative-ai/docs/use-provisioned-throughput
				"X-Vertex-AI-LLM-Request-Type": []string{"shared"},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash"
	contents := genai.Text("How does AI work?")

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)

	// Example response:
	// Artificial Intelligence (AI) isn't magic, nor is it a single "thing." Instead, it's a broad field of computer science focused on creating machines that can perform tasks that typically require human intelligence.
	// .....
	// In Summary:
	// ...

	return nil
}

// [END googlegenaisdk_provisionedthroughput_with_txt]
