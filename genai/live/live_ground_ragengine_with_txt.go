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

// [START googlegenaisdk_live_ground_ragengine_with_txt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateLiveRAGWithText demonstrates how to use the Live API with a Vertex RAG Store.
// It sends a question to the model and retrieves grounded answers from the configured memory corpus.
func generateLiveRAGWithText(w io.Writer, memoryCorpus string) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.0-flash-live-preview-04-09"

	// Configure Vertex RAG store
	ragStore := &genai.VertexRAGStore{
		RAGResources: []*genai.VertexRAGStoreRAGResource{
			{
				RAGCorpus: memoryCorpus, // Define the memory corpus where context is stored or retrieved
			},
		},
	}

	config := &genai.LiveConnectConfig{
		ResponseModalities: []genai.Modality{genai.ModalityText},
		Tools: []*genai.Tool{
			{
				Retrieval: &genai.Retrieval{
					VertexRAGStore: ragStore,
				},
			},
		},
	}

	session, err := client.Live.Connect(ctx, modelName, config)
	if err != nil {
		return fmt.Errorf("failed to connect live: %w", err)
	}
	defer session.Close()

	inputText := "What are the newest Gemini models?"
	fmt.Fprintf(w, "> %s\n\n", inputText)

	// Send the user message
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

		if chunk.ServerContent == nil {
			continue
		}

		// If the server provided a model turn, iterate its parts for text.
		if chunk.ServerContent.ModelTurn != nil {
			for _, part := range chunk.ServerContent.ModelTurn.Parts {
				if part == nil {
					continue
				}
				if part.Text != "" {
					response += part.Text
				}
			}
		}
	}

	fmt.Fprintln(w, response)

	// Example output:
	// > What are the newest Gemini models?
	// In December 2023, Google launched Gemini, their most capable and general model...
	return nil
}

// [END googlegenaisdk_live_ground_ragengine_with_txt]
