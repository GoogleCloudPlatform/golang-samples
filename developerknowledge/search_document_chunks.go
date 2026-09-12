// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package developerknowledge

// [START developerknowledge_search_document_chunks]
import (
	"context"
	"fmt"
	"io"

	developerknowledge "cloud.google.com/go/developerknowledge/apiv1"
	developerknowledgepb "cloud.google.com/go/developerknowledge/apiv1/developerknowledgepb"
	"google.golang.org/api/iterator"
)

// searchDocumentChunks searches developer documentation chunks for a given query.
func searchDocumentChunks(w io.Writer, query string, pageSize int32) ([]*developerknowledgepb.DocumentChunk, error) {
	ctx := context.Background()

	client, err := developerknowledge.NewDeveloperKnowledgeClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("developerknowledge.NewDeveloperKnowledgeClient: %w", err)
	}
	defer client.Close()

	req := &developerknowledgepb.SearchDocumentChunksRequest{
		Query:    query,
		PageSize: pageSize,
	}

	var results []*developerknowledgepb.DocumentChunk
	it := client.SearchDocumentChunks(ctx, req)
	for {
		chunk, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("SearchDocumentChunks: %w", err)
		}
		results = append(results, chunk)
		fmt.Fprintf(w, "Parent Document: %s\n", chunk.GetParent())
		fmt.Fprintf(w, "Chunk ID: %s\n", chunk.GetId())
		fmt.Fprintf(w, "Content: %s\n\n", chunk.GetContent())

		if pageSize > 0 && len(results) >= int(pageSize) {
			break
		}
	}

	return results, nil
}

// [END developerknowledge_search_document_chunks]
