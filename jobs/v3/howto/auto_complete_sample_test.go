// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package howto

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestJobTitleAutoComplete(t *testing.T) {
	tc := testutil.SystemTest(t)

	buf := &bytes.Buffer{}
	if _, err := jobTitleAutoComplete(buf, tc.ProjectID, "", "SWE"); err != nil {
		t.Fatalf("jobTitleAutoComplete: %v", err)
	}
	want := "Auto complete results"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("jobTitleAutoComplete got %q, want %q", got, want)
	}

	buf.Reset()
	if _, err := defaultAutoComplete(buf, tc.ProjectID, "", "SWE"); err != nil {
		t.Fatalf("defaultAutoComplete: %v", err)
	}
	want = "Auto complete results"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("defaultAutoComplete got %q, want %q", got, want)
	}
}

func TestDefaultAutoComplete(t *testing.T) {
	tc := testutil.SystemTest(t)

	buf := &bytes.Buffer{}
	if _, err := defaultAutoComplete(buf, tc.ProjectID, "", "SWE"); err != nil {
		t.Fatalf("defaultAutoComplete: %v", err)
	}
	want := "Auto complete results"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("defaultAutoComplete got %q, want %q", got, want)
	}
}
