// Copyright 2019 Google LLC
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

package log

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLogEntries(t *testing.T) {
	// TODO: Use testutil to get the project.
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if projectID == "" {
		t.Skip("Missing GOLANG_SAMPLES_PROJECT_ID.")
	}

	buf := new(bytes.Buffer)
	if err := logEntries(buf, projectID); err != nil {
		t.Fatalf("logEntries: %v", err)
	}
	want := "Entries:"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("logEntries got %q, want to contain %q", got, want)
	}
}
