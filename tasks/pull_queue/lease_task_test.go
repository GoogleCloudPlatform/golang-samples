// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"io/ioutil"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestTaskLeaseAndAck(t *testing.T) {
	tc := testutil.SystemTest(t)

	// ProjectID is set from environment variable GOLANG_SAMPLES_PROJECT_ID.
	projectID := tc.ProjectID
	locationID := "us-central1"
	queueID := "my-pull-queue"

	// Guarantee a task will be available in the queue
	// in the event TestTaskCreate is skipped.
	seed, err := taskCreate(ioutil.Discard, projectID, locationID, queueID)
	if err != nil {
		t.Fatalf("failed to ensure a task would be available for lease: %v", err)
	}

	// Test task leasing.
	task, err := taskLease(ioutil.Discard, projectID, locationID, queueID)
	if err != nil {
		t.Fatalf("failed to lease a task: %v", err)
	}
	if task == nil {
		t.Fatalf("no task available to lease: %v", err)
	}

	// Note cross-test data poisoning through concurrent queue usage or pre-created tasks.
	if task.GetName() != seed.GetName() {
		t.Logf("Task used for lease testing was not created by test setup: (seeded: %s, leased: %s)", seed.GetName(), task.GetName())
	}

	// Acknowledge our leased task.
	if err = taskAck(ioutil.Discard, task); err != nil {
		t.Fatalf("failed to acknowledge task: %v", err)
	}
}
