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

// [START googlegenaisdk_contentcache_list]
import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"google.golang.org/genai"
)

// listContentCache shows how to retrieve details about cached content.
func listContentCache(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// Retrieve cached content metadata
	cache, err := client.Caches.List(ctx, &genai.ListCachedContentsConfig{
		HTTPOptions: &genai.HTTPOptions{
			Headers:    http.Header{"X-Custom-Header": []string{"example"}},
			APIVersion: "v1",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get content cache: %w", err)
	}

	// Print basic info about the cached content
	fmt.Fprintf(w, "Cache name: %s\n", cache.Name)
	fmt.Fprintf(w, "Display name: %s\n", cache.Items[0].DisplayName)
	fmt.Fprintf(w, "Model: %s\n", cache.Items[0].Model)
	fmt.Fprintf(w, "Create time: %s\n", cache.Items[0].CreateTime.Format(time.RFC3339))
	fmt.Fprintf(w, "Update time: %s\n", cache.Items[0].UpdateTime.Format(time.RFC3339))
	fmt.Fprintf(w, "Expire time: %s (in %s)\n", cache.Items[0].ExpireTime.Format(time.RFC3339), time.Until(cache.Items[0].ExpireTime).Round(time.Second))

	if cache.Items[0].UsageMetadata != nil {
		fmt.Fprintf(w, "Usage metadata: %+v\n", cache.Items[0].UsageMetadata)
	}

	// Example response:
	// Cache name: projects/111111111111/locations/us-central1/cachedContents/1234567890123456789
	// Display name: product_recommendations_prompt
	// Model: models/gemini-2.5-flash
	// Create time: 2025-04-08T02:15:23Z
	// Update time: 2025-04-08T03:05:11Z
	// Expire time: 2025-04-20T03:05:11Z (in 167h59m59s)
	// Usage metadata: &{AudioDurationSeconds:0 ImageCount:167 TextCount:153 TotalTokenCount:43124 VideoDurationSeconds:0}

	return nil
}

// [END googlegenaisdk_contentcache_list]
