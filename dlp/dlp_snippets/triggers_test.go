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

package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

func TestTriggersSamples(t *testing.T) {
	testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	createTrigger(buf, client, projectID, dlppb.Likelihood_POSSIBLE, 0, "my-trigger", "My Trigger", "Test trigger", "my-bucket", true, 10, nil)
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
