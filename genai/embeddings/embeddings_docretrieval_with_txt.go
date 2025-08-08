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

// Package embeddings shows examples of how Gemini-embedding model can use embedding.
package embeddings

// [START googlegenaisdk_embeddings_docretrieval_with_txt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateEmbedContentWithText shows how to embed content with text.
func generateEmbedContentWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	outputDimensionality := int32(3072)
	config := &genai.EmbedContentConfig{
		TaskType:             "RETRIEVAL_DOCUMENT",  //optional
		Title:                "Driver's License",    //optional
		OutputDimensionality: &outputDimensionality, //optional
	}

	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{
					Text: "How do I get a driver's license/learner's permit?",
				},
				{
					Text: "How long is my driver's license valid for?",
				},
				{
					Text: "Driver's knowledge test study guide",
				},
			},
			Role: "user",
		},
	}

	modelName := "gemini-embedding-001"

	resp, err := client.Models.EmbedContent(ctx, modelName, contents, config)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp

	fmt.Fprintln(w, respText)

	// Example response:
	// embeddings=[ContentEmbedding(values=[-0.06302902102470398, 0.00928034819662571, 0.014716853387653828, -0.028747491538524628, ... ],
	// statistics=ContentEmbeddingStatistics(truncated=False, token_count=13.0))]
	// metadata=EmbedContentMetadata(billable_character_count=112)

	return nil
}

// [END googlegenaisdk_embeddings_docretrieval_with_txt]
