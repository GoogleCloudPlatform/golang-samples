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

	developerknowledge "cloud.google.com/go/developerknowledge/apiv1"
	developerknowledgepb "cloud.google.com/go/developerknowledge/apiv1/developerknowledgepb"
)

// batchGetDocuments retrieves multiple developer documentation pages in a single request.
func batchGetDocuments(w io.Writer, names []string) (*developerknowledgepb.BatchGetDocumentsResponse, error) {
	ctx := context.Background()

	client, err := developerknowledge.NewDeveloperKnowledgeClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("developerknowledge.NewDeveloperKnowledgeClient: %w", err)
	}
	defer client.Close()

	req := &developerknowledgepb.BatchGetDocumentsRequest{
		Names: names,
	}

	resp, err := client.BatchGetDocuments(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("BatchGetDocuments: %w", err)
	}

	for _, doc := range resp.GetDocuments() {
		fmt.Fprintf(w, "Title: %s\n", doc.GetTitle())
		fmt.Fprintf(w, "URI: %s\n", doc.GetUri())
		fmt.Fprintf(w, "Content Length: %d bytes\n\n", doc.GetContentLengthBytes())
	}

	return resp, nil
}

// [END developerknowledge_batch_get_documents]
