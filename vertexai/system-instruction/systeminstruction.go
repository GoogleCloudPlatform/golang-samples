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

// systeminstruction shows an example of providing the model with additional context to understand the task
// and get customized responses.
// For developers, product-level behavior can be specified in system instructions (also known as system
// prompts), separate from prompts provided by end users.
package systeminstruction

// [START generativeaionvertexai_gemini_system_instruction]
import (
	"context"
	"errors"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

// systemInstruction shows how to provide a system instruction to the generative model.
func systemInstruction(w io.Writer, projectID, location, modelName string) error {
	// location := "us-central1"
	// modelName := "gemini-1.5-flash-001"

	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	// The System Instruction is set at model creation
	model := client.GenerativeModel(modelName)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(`
			You are a helpful language translator.
			Your mission is to translate text in English to French.
		`)},
	}

	res, err := model.GenerateContent(ctx, genai.Text(`
		User input: I like bagels.
		Answer:
	`))
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

// [END generativeaionvertexai_gemini_system_instruction]
