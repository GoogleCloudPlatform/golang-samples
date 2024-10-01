// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// contextcaching shows an example of caching the tokens of a multimodal PDF prompt
package contextcaching

// [START generativeaionvertexai_gemini_create_context_cache]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/vertexai/genai"
)

// createContextCache shows how to create a cached content, and returns its name.
func createContextCache(w io.Writer, projectID, location, modelName string) (string, error) {
	// location := "us-central1"
	// modelName := "gemini-1.5-pro-001"
	ctx := context.Background()

	systemInstruction := `
    	You are an expert researcher. You always stick to the facts in the sources provided, and never make up new facts.
    	Now look at these research papers, and answer the following questions.
    `

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return "", fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	// These PDF are viewable at
	//   https://storage.googleapis.com/cloud-samples-data/generative-ai/pdf/2312.11805v3.pdf
	//   https://storage.googleapis.com/cloud-samples-data/generative-ai/pdf/2403.05530.pdf

	part1 := genai.FileData{
		MIMEType: "application/pdf",
		FileURI:  "gs://cloud-samples-data/generative-ai/pdf/2312.11805v3.pdf",
	}

	part2 := genai.FileData{
		MIMEType: "application/pdf",
		FileURI:  "gs://cloud-samples-data/generative-ai/pdf/2403.05530.pdf",
	}

	content := &genai.CachedContent{
		Model: modelName,
		SystemInstruction: &genai.Content{
			Parts: []genai.Part{genai.Text(systemInstruction)},
		},
		Expiration: genai.ExpireTimeOrTTL{TTL: 60 * time.Minute},
		Contents: []*genai.Content{
			{
				Role:  "user",
				Parts: []genai.Part{part1, part2},
			},
		},
	}

	result, err := client.CreateCachedContent(ctx, content)
	if err != nil {
		return "", fmt.Errorf("CreateCachedContent: %w", err)
	}
	fmt.Fprint(w, result.Name)
	return result.Name, nil
}

// [END generativeaionvertexai_gemini_create_context_cache]
// FIXME: vburlaka