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

// generateWithLocalImgAndCodeExec shows how to combine a local image, a text prompt,
// and enable the code execution tool in a request.
func generateWithLocalImgAndCodeExec(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// Read local image
	imgBytes, err := os.ReadFile("test_data/640px-Monty_open_door.svg.png")
	if err != nil {
		return fmt.Errorf("failed to read image: %w", err)
	}

	// Define the prompt
	prompt := `
		Run a simulation of the Monty Hall Problem with 1,000 trials.
		Here's how this works as a reminder. In the Monty Hall Problem, you're on a game
		show with three doors. Behind one is a car, and behind the others are goats. You
		pick a door. The host, who knows what's behind the doors, opens a different door
		to reveal a goat. Should you switch to the remaining unopened door?
		The answer has always been a little difficult for me to understand when people
		solve it with math - so please run a simulation with Python to show me what the
		best strategy is.
		Thank you!
	`

	// Enable the code execution tool
	tools := []*genai.Tool{
		{CodeExecution: &genai.ToolCodeExecution{}},
	}

	// Build contents with image + text
	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{InlineData: &genai.Blob{
					MIMEType: "image/png",
					Data:     imgBytes,
				}},
				{Text: prompt},
			},
		},
	}

	// Call the model
	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", contents, &genai.GenerateContentConfig{
		Tools:       tools,
		Temperature: genai.Ptr(float32(0.0)),
	})
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	// Print result
	fmt.Fprintln(w, "# Code:")
	fmt.Fprintln(w, resp.ExecutableCode())
	fmt.Fprintln(w, "# Outcome:")
	fmt.Fprintln(w, resp.CodeExecutionResult())

	// Example output:
	// # Code:
	// import random
	//
	// def run_monty_hall_trial(strategy):
	//      """
	//    Runs a single trial of the Monty Hall problem.
	//
	//    Args:
	//        strategy (str): 'stick' or 'switch'

	return nil
}

// [END googlegenaisdk_tools_code_exec_with_txt_local_img]
