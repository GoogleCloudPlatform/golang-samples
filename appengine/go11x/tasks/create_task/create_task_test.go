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
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"testing"
)

func TestCreateTask(t *testing.T) {
	tc := testutil.SystemTest(t)
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
		_, err := createTask(tc.ProjectID, locationID, queueID, test.message)
		if err != nil {
			t.Errorf("CreateTask(%s): %v", test.name, err)
		}
	}
}
