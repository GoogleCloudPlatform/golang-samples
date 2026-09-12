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

import (
	"bytes"
	"strings"
	"testing"
)

// documentNames uses search to discover the names of real documents so tests
// don't depend on hard-coded document paths, which may move over time.
func documentNames(t *testing.T, query string, pageSize int32) []string {
	t.Helper()
	var buf bytes.Buffer
	chunks, err := searchDocumentChunks(&buf, query, pageSize)
	if err != nil {
		t.Fatalf("searchDocumentChunks(%q): %v", query, err)
	}
	var names []string
	seen := map[string]bool{}
	for _, chunk := range chunks {
		name := chunk.GetParent()
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		names = append(names, name)
	}
	if len(names) == 0 {
		t.Fatalf("searchDocumentChunks(%q): got 0 parent documents, want at least 1", query)
	}
	return names
}

func TestSearchDocumentChunks(t *testing.T) {
	var buf bytes.Buffer
	results, err := searchDocumentChunks(&buf, "Cloud Storage bucket creation", 3)
	if err != nil {
		t.Fatalf("searchDocumentChunks: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("searchDocumentChunks: got 0 results, want non-empty")
	}
	got := buf.String()
	if !strings.Contains(got, "Parent Document: documents/") {
		t.Errorf("searchDocumentChunks output missing parent document prefix, got: %s", got)
	}
}

func TestGetDocument(t *testing.T) {
	var buf bytes.Buffer
	name := documentNames(t, "Cloud Storage bucket creation", 3)[0]
	doc, err := getDocument(&buf, name)
	if err != nil {
		t.Fatalf("getDocument: %v", err)
	}
	if doc.GetName() != name {
		t.Errorf("getDocument: got name %q, want %q", doc.GetName(), name)
	}
	if len(doc.GetTitle()) == 0 {
		t.Errorf("getDocument(%q): got empty title, want non-empty", name)
	}
}

func TestBatchGetDocuments(t *testing.T) {
	var buf bytes.Buffer
	names := documentNames(t, "Cloud Storage buckets", 10)
	if len(names) > 2 {
		names = names[:2]
	}
	resp, err := batchGetDocuments(&buf, names)
	if err != nil {
		t.Fatalf("batchGetDocuments: %v", err)
	}
	if got, want := len(resp.GetDocuments()), len(names); got != want {
		t.Fatalf("batchGetDocuments: got %d documents, want %d", got, want)
	}
}

func TestAnswerQuery(t *testing.T) {
	var buf bytes.Buffer
	resp, err := answerQuery(&buf, "How to create a Cloud Storage bucket")
	if err != nil {
		t.Fatalf("answerQuery: %v", err)
	}
	if resp.GetAnswer() == nil || len(resp.GetAnswer().GetAnswerText()) == 0 {
		t.Errorf("answerQuery: got empty answer text, want non-empty")
	}
}
