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

// Package tools shows examples of various tools that Gemini model can use to generate text.
package tools

// [START googlegenaisdk_tools_code_exec_with_txt]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateWithCodeExec shows how to generate text using the code execution tool.
func generateWithCodeExec(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	prompt := "Calculate 20th fibonacci number. Then find the nearest palindrome to it."
	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: prompt},
		},
			Role: "user"},
	}
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{CodeExecution: &genai.ToolCodeExecution{}},
		},
		Temperature: genai.Ptr(float32(0.0)),
	}
	modelName := "gemini-2.5-flash"

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	for _, p := range resp.Candidates[0].Content.Parts {
		if p.Text != "" {
			fmt.Fprintf(w, "Gemini: %s", p.Text)
		}
		if p.ExecutableCode != nil {
			fmt.Fprintf(w, "Language: %s\n%s\n", p.ExecutableCode.Language, p.ExecutableCode.Code)
		}
		if p.CodeExecutionResult != nil {
			fmt.Fprintf(w, "Outcome: %s\n%s\n", p.CodeExecutionResult.Outcome, p.CodeExecutionResult.Output)
		}
	}

	// Example response:
	// Gemini: Okay, I can do that. First, I'll calculate the 20th Fibonacci number. Then, I need ...
	//
	// Language: PYTHON
	//
	// def fibonacci(n):
	//    ...
	//
	// fib_20 = fibonacci(20)
	// print(f'{fib_20=}')
	//
	// Outcome: OUTCOME_OK
	// fib_20=6765
	//
	// Now that I have the 20th Fibonacci number (6765), I need to find the nearest palindrome. ...
	// ...

	return nil
}

// [END googlegenaisdk_tools_code_exec_with_txt]
