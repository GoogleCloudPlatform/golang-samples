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

	initialContexts, initialErr := listContexts(projectID, sessionID)

	if initialErr != nil {
		t.Error("Unsuccessful initial listContexts")
	}

	var err error
	for _, contextID := range contextIDs {
		err = createContext(projectID, sessionID, contextID)
		if err != nil {
			t.Errorf("Unsuccessful context creation: %s/contexts/%s", parent, contextID)
		}
	}

	intermediateContexts, intermediateErr := listContexts(projectID, sessionID)

	if intermediateErr != nil {
		t.Error("Unsuccessful intermediate listContexts")
	}

	if len(intermediateContexts) != len(initialContexts)+len(contextIDs) {
		t.Errorf("len(intermediateContexts) = %d; want %d", len(intermediateContexts), len(initialContexts)+len(contextIDs))
	}

	for _, contextID := range contextIDs {
		err = deleteContext(projectID, sessionID, contextID)
		if err != nil {
			t.Errorf("Unsuccessful context deletion %s/context/%s", parent, contextID)
		}
	}

	finalContexts, finalErr := listContexts(projectID, sessionID)

	if finalErr != nil {
		t.Error("Unsuccessful final listContexts")
	}

	if len(finalContexts) != len(initialContexts) {
		t.Errorf("Actual len(finalContexts) = %d; want %d", len(finalContexts), len(initialContexts))
	}
}
