// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package tasks

import (
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreateTask(t *testing.T) {
	tc := testutil.SystemTest(t)

	// CreateTask pulls the project ID from the environment.
	os.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)
	locationID := "us-central1"
	queueID := "my-appengine-queue"

	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "Message",
			message: "task details for handler processing",
		},
		{
			name:    "No Message",
			message: "",
		},
	}

	for _, test := range tests {
		_, err := CreateTask(locationID, queueID, test.message)
		if err != nil {
			t.Errorf("CreateTask(%s): %v", test.name, err)
		}
	}
}
