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

// [START googlegenaisdk_tools_code_exec_with_txt_local_img]
import (
	"context"
	"fmt"
	"io"
	"os"

	genai "google.golang.org/genai"
)

// generateWithCodeExecAndImg shows how to generate text using the code execution tool and a local image.
func generateWithCodeExecAndImg(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// Image source:
	//   https://upload.wikimedia.org/wikipedia/commons/thumb/3/3f/Monty_open_door.svg/640px-Monty_open_door.svg.png
	imagePath := filepath.Join(getMedia(), "640px-Monty_open_door.svg.png")
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	prompt := `
Run a simulation of the Monty Hall Problem with 1,000 trials.
Here's how this works as a reminder. In the Monty Hall Problem, you're on a game
show with three doors. Behind one is a car, and behind the others are goats. You
pick a door. The host, who knows what's behind the doors, opens a different door
to reveal a goat. Should you switch to the remaining unopened door?
The answer has always been a little difficult for me to understand when people
solve it with math - so please run a simulation with Python to show me what the
best strategy is.
Thank you!`
	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: prompt},
			{InlineData: &genai.Blob{
				Data:     imgBytes,
				MIMEType: "image/png",
			}},
		}},
	}
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{CodeExecution: &genai.ToolCodeExecution{}},
		},
		Temperature: genai.Ptr(0.0),
	}
	modelName := "gemini-2.0-flash-001"

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
	// Language: PYTHON
	//
	// import random
	//
	// def monty_hall_simulation(num_trials):
	//   ...
	//
	// # Run the simulation for 1000 trials
	// num_trials = 1000
	// switch_win_percentage, stay_win_percentage = monty_hall_simulation(num_trials)
	// ...
	// Outcome: OUTCOME_OK
	// Switching Doors Win Percentage: 63.80%
	// Staying with Original Door Win Percentage: 36.20%
	//
	// Gemini: The results of the simulation clearly demonstrate that switching doors is the better strategy. ...

	return nil
}

// [END googlegenaisdk_tools_code_exec_with_txt_local_img]
