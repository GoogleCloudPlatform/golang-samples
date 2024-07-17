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

// contextcaching shows an example of caching the tokens of a mulitple PDF prompt
package contextcaching

// [START generativeaionvertexai_gemini_update_context_cache]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/vertexai/genai"
)

// updateContextCache shows how to update the expiration time of a cached content, by specifying
// a new TTL (time-to-live duration)
// contentName is the ID of the cached content to update
func updateContextCache(w io.Writer, contentName string, projectID, location string) error {
	// location := "us-central1"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	cachedContent, err := client.GetCachedContent(ctx, contentName)
	if err != nil {
		return fmt.Errorf("GetCachedContent: %w", err)
	}

	update := &genai.CachedContentToUpdate{
		Expiration: &genai.ExpireTimeOrTTL{TTL: 2 * time.Hour},
	}

	_, err = client.UpdateCachedContent(ctx, cachedContent, update)
	fmt.Fprintf(w, "Updated cached content %q", contentName)
	return err
}

// [END generativeaionvertexai_gemini_update_context_cache]
