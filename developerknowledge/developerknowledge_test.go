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

func TestSearchDocumentChunks(t *testing.T) {
	var buf bytes.Buffer
	resp, err := searchDocumentChunks(&buf, "Cloud Storage bucket creation", 3)
	if err != nil {
		t.Fatalf("searchDocumentChunks: %v", err)
	}
	if len(resp.Results) == 0 {
		t.Fatalf("expected non-empty search results")
	}
	got := buf.String()
	if !strings.Contains(got, "Parent Document: documents/") {
		t.Errorf("searchDocumentChunks output missing parent document prefix, got: %s", got)
	}
}

func TestGetDocument(t *testing.T) {
	var buf bytes.Buffer
	name := "documents/docs.cloud.google.com/storage/docs/creating-buckets"
	doc, err := getDocument(&buf, name)
	if err != nil {
		t.Fatalf("getDocument: %v", err)
	}
	if doc.Name != name {
		t.Errorf("got name %q, want %q", doc.Name, name)
	}
	if len(doc.Title) == 0 {
		t.Errorf("expected non-empty title")
	}
}

func TestBatchGetDocuments(t *testing.T) {
	var buf bytes.Buffer
	names := []string{
		"documents/docs.cloud.google.com/storage/docs/creating-buckets",
		"documents/docs.cloud.google.com/storage/docs/deleting-buckets",
	}
	resp, err := batchGetDocuments(&buf, names)
	if err != nil {
		t.Fatalf("batchGetDocuments: %v", err)
	}
	if len(resp.Documents) != 2 {
		t.Fatalf("got %d documents, want 2", len(resp.Documents))
	}
}

func TestAnswerQuery(t *testing.T) {
	var buf bytes.Buffer
	resp, err := answerQuery(&buf, "How to create a Cloud Storage bucket")
	if err != nil {
		t.Fatalf("answerQuery: %v", err)
	}
	if resp.Answer == nil || len(resp.Answer.AnswerText) == 0 {
		t.Errorf("expected non-empty answer text")
	}
}
