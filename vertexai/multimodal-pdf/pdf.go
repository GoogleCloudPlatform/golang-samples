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

// multimodalpdf shows an example of understanding a PDF file input
package multimodalpdf

// [START generativeaionvertexai_gemini_pdf]
import (
	"context"
	"errors"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

// pdfPrompt is a sample prompt type consisting of one PDF asset, and a text question.
type pdfPrompt struct {
	// pdfPath is a Google Cloud Storage path starting with "gs://"
	pdfPath string
	// question asked to the model
	question string
}

// generateContentFromPDF generates a response into the provided io.Writer, based upon the PDF
// asset and the question provided in the multimodal prompt.
func generateContentFromPDF(w io.Writer, prompt pdfPrompt, projectID, location, modelName string) error {
	// prompt := pdfPrompt{
	// 	pdfPath: "gs://cloud-samples-data/generative-ai/pdf/2403.05530.pdf",
	// 	question: `
	// 		You are a very professional document summarization specialist.
	// 		Please summarize the given document.
	// 	`,
	// }
	// location := "us-central1"
	// modelName := "gemini-1.5-pro-preview-0409"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	part := genai.FileData{
		MIMEType: "application/pdf",
		FileURI:  prompt.pdfPath,
	}

	res, err := model.GenerateContent(ctx, part, genai.Text(prompt.question))
	if err != nil {
		return fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 ||
		len(res.Candidates[0].Content.Parts) == 0 {
		return errors.New("empty response from model")
	}

	fmt.Fprintf(w, "generated response: %s\n", res.Candidates[0].Content.Parts[0])
	return nil
}

// [END generativeaionvertexai_gemini_pdf]
