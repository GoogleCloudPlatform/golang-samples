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

// Package live shows how to use the GenAI SDK to generate text with live resources.
package live

// [START googlegenaisdk_live_ground_googsearch_with_txt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateGroundSearchWithTxt demonstrates using a live Gemini model with Google Search grounded responses.
func generateGroundSearchWithTxt(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.0-flash-live-preview-04-09"

	config := &genai.LiveConnectConfig{
		ResponseModalities: []genai.Modality{genai.ModalityText},
		Tools: []*genai.Tool{
			{GoogleSearch: &genai.GoogleSearch{}},
		},
	}

	session, err := client.Live.Connect(ctx, modelName, config)
	if err != nil {
		return fmt.Errorf("failed to connect live session: %w", err)
	}
	defer session.Close()

	textInput := "When did the last Brazil vs. Argentina soccer match happen?"

	// Send user input
	userContent := &genai.Content{
		Role: "user",
		Parts: []*genai.Part{
			{Text: textInput},
		},
	}
	if err := session.SendClientContent(genai.LiveClientContentInput{
		Turns: []*genai.Content{userContent},
	}); err != nil {
		return fmt.Errorf("failed to send client content: %w", err)
	}

	var response string

	// Receive streaming responses
	for {
		chunk, err := session.Receive()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error receiving stream: %w", err)
		}

		// Server setup ready
		if chunk.SetupComplete != nil {
			fmt.Println("Server setup complete")
		}

		// Handle the main model output
		if chunk.ServerContent != nil {
			if chunk.ServerContent.ModelTurn != nil {
				for _, part := range chunk.ServerContent.ModelTurn.Parts {
					if part == nil {
						continue
					}
					if part.Text != "" {
						fmt.Print(part.Text)
						response += part.Text
					}
					if part.ExecutableCode != nil {
						fmt.Println("\n[Executable code]:")
						fmt.Println(part.ExecutableCode.Code)
					}
					if part.CodeExecutionResult != nil {
						fmt.Println("\n[Execution result]:")
						fmt.Println(part.CodeExecutionResult.Output)
					}
				}
			}

			if chunk.ServerContent.GenerationComplete {
				fmt.Println("\n[Generation complete]")
			}

			if chunk.ServerContent.Interrupted {
				fmt.Println("\n[Generation interrupted]")
			}
		}

		if chunk.ToolCall != nil {
			fmt.Println("[ToolCall received]", chunk.ToolCall)
		}

		if chunk.GoAway != nil {
			fmt.Println("[Server requested session end]")
			break
		}
	}

	fmt.Fprintln(w, response)

	// Example output:
	// > When did the last Brazil vs. Argentina soccer match happen?
	// The most recent match between Argentina and Brazil took place on March 25, 2025, as part of the 2026 World Cup qualifiers. Argentina won 4-1.

	return nil
}

// [END googlegenaisdk_live_ground_googsearch_with_txt]
