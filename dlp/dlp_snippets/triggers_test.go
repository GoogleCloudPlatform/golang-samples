// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"strings"
	"testing"

	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

func TestTriggersSamples(t *testing.T) {
	buf := new(bytes.Buffer)
	createTrigger(buf, client, projectID, dlppb.Likelihood_POSSIBLE, 0, "my-trigger", "My Trigger", "Test trigger", "my-bucket", 10, nil)
	if got := buf.String(); !strings.Contains(got, "Successfully created trigger") {
		t.Fatalf("failed to createTrigger: %s", got)
	}
	buf.Reset()
	fullID := "projects/" + projectID + "/jobTriggers/my-trigger"
	listTriggers(buf, client, projectID)
	if got := buf.String(); !strings.Contains(got, fullID) {
		t.Fatalf("failed to list newly created trigger (%s): %q", fullID, got)
	}
	buf.Reset()
	deleteTrigger(buf, client, fullID)
	if got := buf.String(); !strings.Contains(got, "Successfully deleted trigger") {
		t.Fatalf("failed to delete trigger")
	}
}
