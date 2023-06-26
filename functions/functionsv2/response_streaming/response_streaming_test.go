// Copyright 2023 Google LLC
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

package responsestreaming

import (
	"context"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
)

func TestResponseStreaming(t *testing.T) {
	ctx := context.Background()
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if projectID == "" {
		t.Skip("GOLANG_SAMPLES_PROJECT_ID is unset")
	}
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}

	rows, err := query(ctx, client)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	streamResults(w, rows)

	want := true
	if got := strings.Contains(w.Body.String(), "Successfully flushed row"); got != want {
		t.Errorf("Response contains Successfully flushed row: %t, want %t", got, want)
	}
}
