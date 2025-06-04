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

// [START googlegenaisdk_contentcache_use_with_txt]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// useContentCacheWithTxt shows how to use content cache to generate text content.
func useContentCacheWithTxt(w io.Writer, cacheName string) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1beta1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	resp, err := client.Models.GenerateContent(ctx,
		"gemini-1.5-pro-002",
		genai.Text("Summarize the pdfs"),
		&genai.GenerateContentConfig{
			CachedContent: cacheName,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to use content cache to generate content: %w", err)
	}

	respText, err := resp.Text()
	if err != nil {
		return fmt.Errorf("failed to convert model response to text: %w", err)
	}
	fmt.Fprintln(w, respText)

	// Example response:
	// The provided research paper introduces Gemini 1.5 Pro, a multimodal model capable of recalling
	// and reasoning over information from very long contexts (up to 10 million tokens).  Key findings include:
	//
	// * **Long Context Performance:**
	// ...

	return nil
}

// [END googlegenaisdk_contentcache_use_with_txt]
