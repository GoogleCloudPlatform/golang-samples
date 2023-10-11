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

package trigger

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

func TestListTrigger(t *testing.T) {
	tc := testutil.SystemTest(t)

	u := uuid.New().String()[:8]

	triggerID := fmt.Sprintf("golang-list-trigger-test-%s", u)

	var buf bytes.Buffer
	if err := createTriggerForTest(t, tc.ProjectID, triggerID, "My Trigger", "Test trigger", "my-bucket", nil); err != nil {
		t.Fatalf("createTrigger: %v", err)
	}

	fullID := fmt.Sprintf("projects/" + tc.ProjectID + "/locations/global/jobTriggers/" + triggerID)

	if err := listTriggers(&buf, tc.ProjectID); err != nil {
		t.Errorf("listTriggers: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, fullID) {
		t.Errorf("listTriggers got\n----\n%v\n----\nWant to contain:\n----\n%v\n----", got, fullID)
	}

	cleanUpTrigger(t, tc.ProjectID, fullID)
}
