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
package token

import (
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreateHTTPTaskWithToken(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := "us-central1"
	queueID := "my-queue"
	url := "https://example.com/task_handler"
	serviceAccountEmail := os.Getenv("GOLANG_SAMPLES_SERVICE_ACCOUNT_EMAIL")
	if serviceAccountEmail == "" {
		t.Skip("GOLANG_SAMPLES_SERVICE_ACCOUNT_EMAIL not set")
	}

	tests := []struct {
		message string
	}{
		{
			message: "task details for handler processing",
		},
		{
			message: "",
		},
	}

	for _, test := range tests {
		_, err := createHTTPTaskWithToken(tc.ProjectID, locationID, queueID, url, serviceAccountEmail, test.message)
		if err != nil {
			t.Errorf("CreateTask(%q): %v", test.message, err)
		}
	}
}
