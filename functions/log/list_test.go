// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
	buf := new(bytes.Buffer)
	if err := logEntries(buf, projectID); err != nil {
		t.Fatalf("logEntries: %v", err)
	}
	want := "Entries:"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("logEntries got %q, want to contain %q", got, want)
	}
}
