// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tool_use

// [START genai_tool_use_code_execution]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// codeExecution generates code for the given text prompt using Code Execution as a Tool.
func codeExecution(w io.Writer) error {
	modelName := "gemini-2.0-flash-exp"
	client, err := genai.NewClient(context.TODO(), &genai.ClientConfig{})
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}

	codeExecTool := genai.Tool{
		CodeExecution: &genai.ToolCodeExecution{},
	}
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{&codeExecTool},
	}

	textpart := genai.Text(`Calculate 20th fibonacci number. Then find the nearest palindrome to it.`)

	result, err := client.Models.GenerateContent(context.TODO(), modelName,
		&genai.ContentParts{textpart}, config)
	if err != nil {
		return fmt.Errorf("GenerateContent: %w", err)
	}

	for _, part := range result.Candidates[0].Content.Parts {
		if part.ExecutableCode != nil {
			fmt.Fprintf(w, "Code (%s):\n%s\n", part.ExecutableCode.Language, part.ExecutableCode.Code)
		}
		if part.CodeExecutionResult != nil {
			fmt.Fprintf(w, "Result (%s):\n %s\n", part.CodeExecutionResult.Outcome, part.CodeExecutionResult.Output)
		}
	}
	return nil
}

// [END genai_tool_use_code_execution]
