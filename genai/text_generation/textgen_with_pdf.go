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

// [START googlegenaisdk_textgen_with_pdf]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateTextWithPDF shows how to generate text using a PDF file input.
func generateTextWithPDF(w io.Writer) error {
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
			{Text: `You are a highly skilled document summarization specialist.
	Your task is to provide a concise executive summary of no more than 300 words.
	Please summarize the given document for a general audience.`},
			{FileData: &genai.FileData{
				FileURI:  "gs://cloud-samples-data/generative-ai/pdf/1706.03762v7.pdf",
				MIMEType: "application/pdf",
			}},
		},
			Role: genai.RoleUser},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)

	// Example response:
	// "Attention Is All You Need" introduces the Transformer,
	// a groundbreaking neural network architecture designed for...
	// ...

	return nil
}

// [END googlegenaisdk_textgen_with_pdf]
