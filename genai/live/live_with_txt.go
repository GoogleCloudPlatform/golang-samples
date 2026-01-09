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

// [START googlegenaisdk_live_with_txt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateLiveWithText demonstrates using a live Gemini model
// that performs live with text and handles responses.
func generateLiveWithText(w io.Writer) error {
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
	}

	// Open a live session
	session, err := client.Live.Connect(ctx, modelName, config)
	if err != nil {
		return fmt.Errorf("failed to connect live: %w", err)
	}
	defer session.Close()

	// Prepare the input message
	inputText := "Hello? Gemini, are you there?"
	fmt.Fprintf(w, "> %s\n\n", inputText)

	// Send text content to the model
	err = session.SendClientContent(genai.LiveClientContentInput{
		Turns: []*genai.Content{
			{
				Role: genai.RoleUser,
				Parts: []*genai.Part{
					{Text: inputText},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send content: %w", err)
	}

	// Stream the response
	var response string
	for {
		chunk, err := session.Receive()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error receiving response: %w", err)
		}

		if chunk.ServerContent != nil && chunk.ServerContent.ModelTurn != nil {
			for _, part := range chunk.ServerContent.ModelTurn.Parts {
				if part.Text != "" {
					response += part.Text
				}
			}
		}
	}

	// Example output:
	//  >  Hello? Gemini are you there?
	//  Yes, I'm here. What would you like to talk about?
	fmt.Fprintln(w, response)
	return nil
}

// [END googlegenaisdk_live_with_txt]
