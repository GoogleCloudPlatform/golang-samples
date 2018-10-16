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

var projectID string

func assertContains(t *testing.T, out string, sub string) {
	t.Helper()
	if !strings.Contains(out, sub) {
		t.Errorf("got output %q; want it to contain %q", out, sub)
	}
}

func TestMain(m *testing.M) {
	if c, ok := testutil.ContextMain(m); ok {
		projectID = c.ProjectID
	}
	os.Exit(m.Run())
}

func TestInspectString(t *testing.T) {
	testutil.SystemTest(t)
	// Capture snippet log
	var b bytes.Buffer
	log.SetOutput(&b)
	// Run Snippet
	inspectString(projectID, "I'm Gary and my email is gary@example.com")
	log.SetOutput(os.Stderr)
	//Verify snippet output
	output := b.String()
	assertContains(t, output, "Info type: EMAIL_ADDRESS")
}

func TestInspectFile(t *testing.T) {
	testutil.SystemTest(t)
	// Capture snippet log
	var b bytes.Buffer
	log.SetOutput(&b)
	// Run Snippet
	inspectFile(projectID, "testdata/test.png", "IMAGE")
	log.SetOutput(os.Stderr)
	//Verify snippet output
	output := b.String()
	assertContains(t, output, "Info type: PHONE_NUMBER")
	assertContains(t, output, "Info type: EMAIL_ADDRESS")
}
