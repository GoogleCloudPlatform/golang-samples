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

func TestCreateTrigger(t *testing.T) {
	tc := testutil.SystemTest(t)

	u := uuid.New().String()[:8]

	triggerID := fmt.Sprintf("golang-create-trigger-test-%s", u)

	var buf bytes.Buffer
	if err := createTrigger(&buf, tc.ProjectID, triggerID, "My Trigger", "Test trigger", "my-bucket", nil); err != nil {
		t.Fatalf("createTrigger: %v", err)
	}
	if got, want := buf.String(), "Successfully created trigger"; !strings.Contains(got, want) {
		t.Errorf("createTrigger got\n----\n%v\n----\nWant to contain:\n----\n%v\n----", got, want)
	}

	fullID := fmt.Sprintf("projects/%s/locations/global/jobTriggers/%s", tc.ProjectID, triggerID)
	cleanUpTrigger(t, tc.ProjectID, fullID)
}
