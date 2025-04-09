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

// [START googlegenaisdk_textgen_code_with_pdf]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithPDF shows how to generate code using a PDF file input.
func generateWithPDF(w io.Writer) error {
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
			{Text: "Convert this python code to use Google Python Style Guide."},
			{FileData: &genai.FileData{
				FileURI:  "https://storage.googleapis.com/cloud-samples-data/generative-ai/text/inefficient_fibonacci_series_python_code.pdf",
				MIMEType: "application/pdf",
			}},
		}},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText, err := resp.Text()
	if err != nil {
		return fmt.Errorf("failed to convert model response to text: %w", err)
	}
	fmt.Fprintln(w, respText)

	// Example response:
	// ```python
	// def fibonacci(n: int) -> list[int]:
	// ...

	return nil
}

// [END googlegenaisdk_textgen_code_with_pdf]
