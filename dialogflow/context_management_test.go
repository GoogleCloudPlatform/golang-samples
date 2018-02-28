// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestContextManagement(t *testing.T) {
	testutil.SystemTest(t)

	projectId := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if projectId == "" {
		log.Fatal("Please pass a project ID using the GOLANG_SAMPLES_PROJECT_ID environment variable")
	}

	sessionId := fmt.Sprintf("golang-samples-test-session-%v", time.Now())

	parent := fmt.Sprintf("projects/%s/agents/sessions/%s", projectId, sessionId)

	contextIds := [...]string{"context-1", "context-2"}

	initialContexts, initialErr := listContexts(projectId, sessionId)

	if initialErr != nil {
		t.Error("Unsuccessful initial listContexts")
	}

	var err error
	for _, contextId := range contextIds {
		err = createContext(projectId, sessionId, contextId)
		if err != nil {
			t.Errorf("Unsuccessful context creation: %s/contexts/%s", parent, contextId)
		}
	}

	intermediateContexts, intermediateErr := listContexts(projectId, sessionId)

	if intermediateErr != nil {
		t.Error("Unsuccessful intermediate listContexts")
	}

	if len(intermediateContexts) != len(initialContexts) + len(contextIds) {
		t.Errorf("Expected intermediateContexts to contain %d more values than initialContexts. Actual len(intermediateContexts): %d, actual len(initialContexts): %d", len(contextIds), len(intermediateContexts), len(initialContexts))
	}

	for _, contextId := range contextIds {
		err = deleteContext(projectId, sessionId, contextId)
		if err != nil {
			t.Errorf("Unsuccessful context deletion %s/context/%s", parent, contextId)
		}
	}

	finalContexts, finalErr := listContexts(projectId, sessionId)

	if finalErr != nil {
		t.Error("Unsuccessful final listContexts")
	}

	if len(finalContexts) != len(initialContexts) {
		t.Errorf("Expected finalContexts to be of same length as initialContexts. Actual len(finalContexts): %d, actual len(initialContexts): %d", len(finalContexts), len(initialContexts))
	}
}
