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

// [START genai_tool_use_code_execution_fileinput]
import (
	"context"
	"fmt"
	"io"
	"net/http"

	genai "google.golang.org/genai"
)

// codeExecution generates code for the given text and image prompt using Code Execution as a Tool.
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

	resp, err := http.Get("https://upload.wikimedia.org/wikipedia/commons/thumb/3/3f/Monty_open_door.svg/640px-Monty_open_door.svg.png")
	if err != nil {
		return fmt.Errorf("error fetching image: %w", err)
	}
	defer resp.Body.Close()
	imagebytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %w", err)
	}
	imagepart := genai.InlineData{
		Data:     imagebytes,
		MIMEType: "image/png",
	}

	textpart := genai.Text(`
	Run a simulation of the Monty Hall Problem with 1,000 trials.
	  Here's how this works as a reminder. In the Monty Hall Problem, you're on a game
	  show with three doors. Behind one is a car, and behind the others are goats. You
	  pick a door. The host, who knows what's behind the doors, opens a different door
	  to reveal a goat. Should you switch to the remaining unopened door?
	  The answer has always been a little difficult for me to understand when people
	  solve it with math - so please run a simulation with Python to show me what the
	  best strategy is.
	  Thank you!
	`)

	result, err := client.Models.GenerateContent(context.TODO(), modelName,
		&genai.ContentParts{imagepart, textpart}, config)
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

// [END genai_tool_use_code_execution_fileinput]
