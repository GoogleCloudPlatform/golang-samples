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

// [START generativeaionvertexai_context_caching_list]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/iterator"
)

// listContextCaches retrieves all context caches associated with the specified
// Google Cloud project and region
func listContextCaches(w io.Writer, projectID, location string) error {
	// location := "us-central1"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	cacheList := client.ListCachedContents(ctx)
	// `cacheList` is a standard Google API iterator.
	// See https://pkg.go.dev/google.golang.org/api/iterator#example-package for more details
	for {
		item, err := cacheList.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("listContextCaches failed: %w", err)
		}

		fmt.Fprintf(w, "Cache %q will expire at %v\n", item.Name, item.Expiration.ExpireTime.String())
		// Example response:
		// Cache "projects/.../locations/.../cachedContents/12345678900000000" will expire at 2024-10-25 09:13:58.67004 +0000 UTC
	}

	return nil
}

// [END generativeaionvertexai_context_caching_list]
