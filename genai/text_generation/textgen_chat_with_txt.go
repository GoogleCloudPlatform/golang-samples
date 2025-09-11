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

// Package text_generation shows examples of generating text using the GenAI SDK.
package text_generation

// [START googlegenaisdk_textgen_chat_with_txt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateChatWithText shows how to generate chat using a text prompt.
func generateChatWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}
	modelName := "gemini-2.5-flash"
	history := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: "Hello there"},
			},
		},
		{
			Role: "model",
			Parts: []*genai.Part{
				{Text: "Great to meet you. What would you like to know?"},
			},
		},
	}
	chatSession, err := client.Chats.Create(ctx, modelName, nil, history)
	if err != nil {
		return fmt.Errorf("failed to create genai chat session: %w", err)
	}
	contents := genai.Part{Text: "Tell me a story."}
	resp, err := chatSession.SendMessage(ctx, contents)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	respText := resp.Text()

	fmt.Fprintln(w, respText)
	// Example response:
	// Okay, settle in. Let me tell you a story about a quiet cartographer, but not of lands and seas.
	// ...
	// In the sleepy town of Oakhaven, nestled between the Whispering Hills and the Murmuring River, lived a woman named Elara.
	// ...

	return nil
}

// [END googlegenaisdk_textgen_chat_with_txt]
