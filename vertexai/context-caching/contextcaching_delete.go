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

// [START generativeaionvertexai_gemini_delete_context_cache]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

// deleteContextCache shows how to delete a cached content
// contentName is the ID of the cached content
func deleteContextCache(w io.Writer, contentName string, projectID, location string) error {
	// location := "us-central1"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	err = client.DeleteCachedContent(ctx, contentName)
	if err != nil {
		return fmt.Errorf("DeleteCachedContent: %w", err)
	}
	fmt.Fprintf(w, "Deleted cached content %q", contentName)
	return nil
}

// [END generativeaionvertexai_gemini_delete_context_cache]
