// Copyright 2024 Google LLC
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

package firestore

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestQueryExplainAnalyze(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("Skipping firestore test. Set GOLANG_SAMPLES_FIRESTORE_PROJECT.")
	}

	_, cleanup := setupClientAndCities(t, projectID)
	t.Cleanup(cleanup)

	// Run sample and capture console output
	buf := new(bytes.Buffer)
	if err := queryExplainAnalyze(buf, projectID); err != nil {
		t.Errorf("queryWithExplainAnalyze: %v", err)
	}

	// Compare console outputs
	got := buf.String()
	want := "----- Indexes Used -----\n" +
		"0: &map[properties:(__name__ ASC) query_scope:Collection]\n" +
		"----- Execution Stats -----\n" +
		"&{ResultsReturned:"
	if !strings.Contains(got, want) {
		t.Errorf("%q does not contain %q", got, want)
	}

	want = "index_entries_scanned"
	if !strings.Contains(got, want) {
		t.Errorf("%q does not contain %q", got, want)
	}
}
