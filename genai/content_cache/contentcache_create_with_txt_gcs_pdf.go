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

// Package content_cache shows examples of using content caching with the GenAI SDK.
package content_cache

// [START googlegenaisdk_contentcache_create_with_txt_gcs_pdf]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// createContentCache shows how to create a content cache with an expiration parameter.
func createContentCache(w io.Writer) (string, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1beta1"},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.0-flash-001"

	systemInstruction := "You are an expert researcher. You always stick to the facts " +
		"in the sources provided, and never make up new facts. " +
		"Now look at these research papers, and answer the following questions."

	cacheContents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{FileData: &genai.FileData{
					FileURI:  "gs://cloud-samples-data/generative-ai/pdf/2312.11805v3.pdf",
					MIMEType: "application/pdf",
				}},
				{FileData: &genai.FileData{
					FileURI:  "gs://cloud-samples-data/generative-ai/pdf/2403.05530.pdf",
					MIMEType: "application/pdf",
				}},
			},
			Role: "user",
		},
	}
	config := &genai.CreateCachedContentConfig{
		Contents: cacheContents,
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: systemInstruction},
			},
		},
		DisplayName: "example-cache",
		TTL:         "86400s",
	}

	res, err := client.Caches.Create(ctx, modelName, config)
	if err != nil {
		return "", fmt.Errorf("failed to create content cache: %w", err)
	}

	cachedContent, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal cache info: %w", err)
	}

	// See the documentation: https://pkg.go.dev/google.golang.org/genai#CachedContent
	fmt.Fprintln(w, string(cachedContent))

	// Example response:
	// {
	//   "name": "projects/111111111111/locations/us-central1/cachedContents/1111111111111111111",
	//   "displayName": "example-cache",
	//   "model": "projects/111111111111/locations/us-central1/publishers/google/models/gemini-1.5-pro-002",
	//   "createTime": "2025-02-18T15:05:08.29468Z",
	//   "updateTime": "2025-02-18T15:05:08.29468Z",
	//   "expireTime": "2025-02-19T15:05:08.280828Z",
	//   "usageMetadata": {
	//     "imageCount": 167,
	//     "textCount": 153,
	//     "totalTokenCount": 43125
	//   }
	// }

	return res.Name, nil
}

// [END googlegenaisdk_contentcache_create_with_txt_gcs_pdf]
