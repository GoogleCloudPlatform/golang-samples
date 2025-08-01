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

// [START googlegenaisdk_textgen_chat_stream_with_txt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateChatStreamWithText shows how to generate chat stream using a text prompt.
func generateChatStreamWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.5-flash"

	chatSession, err := client.Chats.Create(ctx, modelName, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create genai chat session: %w", err)
	}

	var streamErr error
	contents := genai.Part{Text: "Why is the sky blue?"}

	stream := chatSession.SendMessageStream(ctx, contents)
	stream(func(resp *genai.GenerateContentResponse, err error) bool {
		if err != nil {
			streamErr = err
			return false
		}
		for _, cand := range resp.Candidates {
			for _, part := range cand.Content.Parts {
				fmt.Fprintln(w, part.Text)
			}
		}
		return true
	})

	// Example response:
	// The
	// sky appears blue due to a phenomenon called **Rayleigh scattering**.
	// Here's a breakdown:
	// ...

	return streamErr
}

// [END googlegenaisdk_textgen_chat_stream_with_txt]
