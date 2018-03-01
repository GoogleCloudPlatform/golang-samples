// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
