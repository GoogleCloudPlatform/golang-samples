// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"io/ioutil"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestTaskCreate(t *testing.T) {
	tc := testutil.SystemTest(t)

	// ProjectID is set from environment variable GOLANG_SAMPLES_PROJECT_ID.
	projectID := tc.ProjectID
	locationID := "us-central1"
	queueID := "my-pull-queue"

	_, err := taskCreate(ioutil.Discard, projectID, locationID, queueID)
	if err != nil {
		t.Fatalf("failed to create new task: %v", err)
	}
}
