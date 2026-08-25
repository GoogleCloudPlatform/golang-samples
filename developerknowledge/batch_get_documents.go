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

// [START developerknowledge_batch_get_documents]
import (
	"context"
	"fmt"
	"io"

	developerknowledge "google.golang.org/api/developerknowledge/v1"
)

// batchGetDocuments retrieves multiple developer documentation pages in a single request.
func batchGetDocuments(
	w io.Writer,
	names []string,
) (*developerknowledge.BatchGetDocumentsResponse, error) {
	ctx := context.Background()

	svc, err := developerknowledge.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("developerknowledge.NewService: %w", err)
	}

	resp, err := svc.Documents.BatchGet().Names(names...).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("BatchGetDocuments: %w", err)
	}

	for _, doc := range resp.Documents {
		fmt.Fprintf(w, "Title: %s\n", doc.Title)
		fmt.Fprintf(w, "URI: %s\n", doc.Uri)
		fmt.Fprintf(w, "Content Length: %d bytes\n\n", doc.ContentLengthBytes)
	}

	return resp, nil
}

// [END developerknowledge_batch_get_documents]
