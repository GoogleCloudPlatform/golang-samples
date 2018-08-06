// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var projectID, locationID, queueID string

func setup(t *testing.T) {
	tc := testutil.SystemTest(t)

	// ProjectID is set from environment variable GOLANG_SAMPLES_PROJECT_ID.
	projectID = tc.ProjectID
        locationID = "us-central1"
    queueID = "my-pull-queue"
}

func TestTaskCreate(t *testing.T) {
	setup(t)

	_, err := taskCreate(projectID, locationID, queueID)
	if err != nil {
		t.Fatalf("failed to create new task: %v", err)
	}
}

func TestTaskLeaseAndAck(t *testing.T) {
	setup(t)

	// Guarantee a task will be available in the queue
	// in the event TestTaskCreate is skipped.
	_, err := taskCreate(projectID, locationID, queueID)
	if err != nil {
		t.Fatalf("failed to ensure a task would be available for lease: %v", err)
	}

	// Test task leasing.
	task, err := taskLease(projectID, locationID, queueID)
	if err != nil {
		t.Fatalf("failed to lease a task: %v", err)
        }
        if task == nil {
		t.Fatalf("no task available to lease: %v", err)
	}

	// Acknowledge our leased task.
	if err = taskAck(task); err != nil {
		t.Fatalf("failed to acknowledge task: %v", err)
	}
}
