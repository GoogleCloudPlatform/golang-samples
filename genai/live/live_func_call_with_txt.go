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

// [START googlegenaisdk_live_func_call_with_txt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateLiveFuncCallWithTxt demonstrates using a live Gemini model
// that performs function calls and handles responses.
func generateLiveFuncCallWithTxt(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelID := "gemini-2.0-flash-live-preview-04-09"

	// Define simple function declarations.
	turnOnLights := &genai.FunctionDeclaration{Name: "turn_on_the_lights"}
	turnOffLights := &genai.FunctionDeclaration{Name: "turn_off_the_lights"}

	config := &genai.LiveConnectConfig{
		ResponseModalities: []genai.Modality{genai.ModalityText},
		Tools: []*genai.Tool{
			{
				FunctionDeclarations: []*genai.FunctionDeclaration{
					turnOnLights,
					turnOffLights,
				},
			},
		},
	}

	session, err := client.Live.Connect(ctx, modelID, config)
	if err != nil {
		return fmt.Errorf("failed to connect live session: %w", err)
	}

	textInput := "Turn on the lights please"
	fmt.Fprintf(w, "> %s\n\n", textInput)

	// Send the user's text as a live content message.
	if err := session.SendClientContent(genai.LiveClientContentInput{
		Turns: []*genai.Content{
			{
				Role: "user",
				Parts: []*genai.Part{
					{Text: textInput},
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to send client content: %w", err)
	}

	var functionResponses []*genai.FunctionResponse

	for {
		chunk, err := session.Receive()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error receiving chunk: %w", err)
		}

		// Handle model-generated content
		if chunk.ServerContent != nil && chunk.ServerContent.ModelTurn != nil {
			for _, part := range chunk.ServerContent.ModelTurn.Parts {
				if part.Text != "" {
					fmt.Fprint(w, part.Text)
				}
			}
		}

		// Handle tool (function) calls
		if chunk.ToolCall != nil {
			for _, fc := range chunk.ToolCall.FunctionCalls {
				functionResponse := &genai.FunctionResponse{
					Name: fc.Name,
					Response: map[string]any{
						"result": "ok",
					},
				}
				functionResponses = append(functionResponses, functionResponse)
				fmt.Fprintln(w, functionResponse.Response["result"])
			}

			if err := session.SendToolResponse(genai.LiveToolResponseInput{
				FunctionResponses: functionResponses,
			}); err != nil {
				return fmt.Errorf("failed to send tool response: %w", err)
			}
		}
	}

	// Example output:
	// >  Turn on the lights please
	// ok

	return nil
}

// [END googlegenaisdk_live_func_call_with_txt]
