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
	"fmt"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestContextManagement(t *testing.T) {
	tc := testutil.SystemTest(t)

	projectID := tc.ProjectID

	sessionID := fmt.Sprintf("golang-samples-test-session-%v", time.Now())

	parent := fmt.Sprintf("projects/%s/agents/sessions/%s", projectID, sessionID)

	contextIDs := [...]string{"context-1", "context-2"}

	initialContexts, err := ListContexts(projectID, sessionID)

	if err != nil {
		t.Error("Unsuccessful initial ListContexts")
	}

	for _, contextID := range contextIDs {
		err = CreateContext(projectID, sessionID, contextID)
		if err != nil {
			t.Errorf("Unsuccessful context creation: %s/contexts/%s", parent, contextID)
		}
	}

	intermediateContexts, err := ListContexts(projectID, sessionID)

	if err != nil {
		t.Error("Unsuccessful intermediate ListContexts")
	}

	if len(intermediateContexts) != len(initialContexts)+len(contextIDs) {
		t.Errorf("len(intermediateContexts) = %d; want %d", len(intermediateContexts), len(initialContexts)+len(contextIDs))
	}

	for _, contextID := range contextIDs {
		err = DeleteContext(projectID, sessionID, contextID)
		if err != nil {
			t.Errorf("Unsuccessful context deletion %s/context/%s", parent, contextID)
		}
	}

	finalContexts, err := ListContexts(projectID, sessionID)

	if err != nil {
		t.Error("Unsuccessful final ListContexts")
	}

	if len(finalContexts) != len(initialContexts) {
		t.Errorf("Actual len(finalContexts) = %d; want %d", len(finalContexts), len(initialContexts))
	}
}
