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

	developerknowledge "google.golang.org/api/developerknowledge/v1"
)

// searchDocumentChunks searches developer documentation chunks for a given query.
func searchDocumentChunks(
	w io.Writer,
	query string,
	pageSize int64,
) (*developerknowledge.SearchDocumentChunksResponse, error) {
	ctx := context.Background()

	svc, err := developerknowledge.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("developerknowledge.NewService: %w", err)
	}

	call := svc.Documents.SearchDocumentChunks().Query(query).PageSize(pageSize).Context(ctx)
	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("SearchDocumentChunks: %w", err)
	}

	for _, chunk := range resp.Results {
		fmt.Fprintf(w, "Parent Document: %s\n", chunk.Parent)
		fmt.Fprintf(w, "Chunk ID: %s\n", chunk.Id)
		fmt.Fprintf(w, "Content: %s\n\n", chunk.Content)
	}

	return resp, nil
}

// [END developerknowledge_search_document_chunks]
