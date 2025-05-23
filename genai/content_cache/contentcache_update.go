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

// [START googlegenaisdk_contentcache_update]
import (
	"context"
	"fmt"
	"io"
	"time"

	genai "google.golang.org/genai"
)

// updateContentCache shows how to update content cache expiration time.
func updateContentCache(w io.Writer, cacheName string) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// Update expire time using TTL
	resp, err := client.Caches.Update(ctx, cacheName, &genai.UpdateCachedContentConfig{
		TTL: int64(36000),
	})
	if err != nil {
		return fmt.Errorf("failed to update content cache exp. time with TTL: %w", err)
	}

	fmt.Fprintf(w, "Cache expires in: %s\n", time.Until(resp.ExpireTime))
	// Example response:
	// Cache expires in: 10h0m0.005875s

	// Update expire time using specific time stamp
	inSevenDays := time.Now().Add(7 * 24 * time.Hour)
	resp, err = client.Caches.Update(ctx, cacheName, &genai.UpdateCachedContentConfig{
		ExpireTime: inSevenDays,
	})
	if err != nil {
		return fmt.Errorf("failed to update content cache expire time: %w", err)
	}

	fmt.Fprintf(w, "Cache expires in: %s\n", time.Until(resp.ExpireTime))
	// Example response:
	// Cache expires in: 167h59m59.80327s

	return nil
}

// [END googlegenaisdk_contentcache_update]
