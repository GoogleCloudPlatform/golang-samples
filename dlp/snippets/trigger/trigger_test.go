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

package trigger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestTriggersSamples(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)

	fullID := "projects/" + tc.ProjectID + "/jobTriggers/my-trigger"

	// Delete the trigger if it already exists since the same ID is used every
	// time.
	if err := listTriggers(buf, tc.ProjectID); err != nil {
		t.Errorf("listTriggers: %v", err)
	}
	if got := buf.String(); strings.Contains(got, fullID) {
		buf.Reset()
		if err := deleteTrigger(buf, fullID); err != nil {
			t.Errorf("deleteTrigger: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, "Successfully deleted trigger") {
			t.Error("failed to delete trigger")
		}
	}

	if err := createTrigger(buf, tc.ProjectID, "my-trigger", "My Trigger", "Test trigger", "my-bucket", nil); err != nil {
		t.Errorf("createTrigger: %v", err)
	}
	if got, want := buf.String(), "Successfully created trigger"; !strings.Contains(got, want) {
		t.Errorf("createTrigger got\n----\n%v\n----\nWant to contain:\n----\n%v\n----", got, want)
	}

	buf.Reset()
	if err := listTriggers(buf, tc.ProjectID); err != nil {
		t.Errorf("listTriggers: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, fullID) {
		t.Errorf("listTriggers got\n----\n%v\n----\nWant to contain:\n----\n%v\n----", got, fullID)
	}

	buf.Reset()
	if err := deleteTrigger(buf, fullID); err != nil {
		t.Errorf("deleteTrigger: %v", err)
	}
	if got, want := buf.String(), "Successfully deleted trigger"; !strings.Contains(got, want) {
		t.Errorf("deleteTrigger got\n----\n%v\n----\nWant to contain:\n----\n%v\n----", got, want)
	}
}
