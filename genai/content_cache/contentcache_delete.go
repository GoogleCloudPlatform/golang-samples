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

// [START googlegenaisdk_contentcache_delete]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// deleteContentCache shows how to delete content cache.
func deleteContentCache(w io.Writer, cacheName string) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1beta1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	_, err = client.Caches.Delete(ctx, cacheName, &genai.DeleteCachedContentConfig{})
	if err != nil {
		return fmt.Errorf("failed to delete content cache: %w", err)
	}

	fmt.Fprintf(w, "Deleted cache %q\n", cacheName)

	// Example response:
	// Deleted cache "projects/111111111111/locations/us-central1/cachedContents/1111111111111111111"

	return nil
}

// [END googlegenaisdk_contentcache_delete]
