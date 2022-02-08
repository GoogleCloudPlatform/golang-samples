// Copyright 2022 Google LLC
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
	"context"
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestAnalyzeIAMPolicyGCS(t *testing.T) {
	tc := testutil.SystemTest(t)
	scope := fmt.Sprintf("projects/%s", tc.ProjectID)
	fullResourceName := fmt.Sprintf("//cloudresourcemanager.googleapis.com/projects/%s", tc.ProjectID)
	// Delete the bucket (if it exists) then recreate it.
	bucketName := fmt.Sprintf("%s-for-assets", tc.ProjectID)
	testutil.CleanBucket(context.Background(), t, tc.ProjectID, bucketName)
	uri := fmt.Sprintf("gs://%s/client_library_obj", bucketName)

	if err := analyzeIAMPolicyGCS(scope, fullResourceName, uri); err != nil {
		t.Errorf("execution failed: %v", err)
	}
}
