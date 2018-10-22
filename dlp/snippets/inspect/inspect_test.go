// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package inspect

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestInspectString(t *testing.T) {
	tc := testutil.SystemTest(t)
	// Run snippet and capture output.
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	err := inspectString(buf, tc.ProjectID, "I'm Gary and my email is gary@example.com")
	log.SetOutput(os.Stderr)

	if err != nil {
		t.Errorf("TestInspectFile: %v", err)
	}
	got := buf.String()
	if want := "Info type: EMAIL_ADDRESS"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInspectFile(t *testing.T) {
	tc := testutil.SystemTest(t)
	// Run snippet and capture output.
	buf := new(bytes.Buffer)
	err := inspectFile(buf, tc.ProjectID, "testdata/test.png", "IMAGE")
	log.SetOutput(os.Stderr)

	if err != nil {
		t.Errorf("TestInspectFile: %v", err)
	}
	got := buf.String()
	if want := "Info type: PHONE_NUMBER"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
	if want := "Info type: EMAIL_ADDRESS"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
