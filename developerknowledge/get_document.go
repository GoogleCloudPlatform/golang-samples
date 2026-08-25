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

// [START developerknowledge_get_document]
import (
	"context"
	"fmt"
	"io"

	developerknowledge "google.golang.org/api/developerknowledge/v1"
)

// getDocument retrieves a single developer documentation page by its resource name.
func getDocument(w io.Writer, name string) (*developerknowledge.Document, error) {
	ctx := context.Background()

	svc, err := developerknowledge.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("developerknowledge.NewService: %w", err)
	}

	doc, err := svc.Documents.Get(name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("GetDocument: %w", err)
	}

	fmt.Fprintf(w, "Title: %s\n", doc.Title)
	fmt.Fprintf(w, "URI: %s\n", doc.Uri)
	fmt.Fprintf(w, "Data Source: %s\n", doc.DataSource)
	fmt.Fprintf(w, "Content Length: %d bytes\n\n", doc.ContentLengthBytes)

	return doc, nil
}

// [END developerknowledge_get_document]
