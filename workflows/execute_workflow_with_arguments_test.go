// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package workflows

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// TestExecuteWorkflowWithArguments tests the
// executeWorkflowWithArguments function and
// evaluates the success by comparing if the
// function's output contains an expected value.
func TestExecuteWorkflowWithArguments(t *testing.T) {
	tc := testutil.SystemTest(t)

	workflowID := testGenerateWorkflowID()
	locationID := "us-central1"

	var buf bytes.Buffer

	// Create the test workflow that will be cleaned up once the test is done.
	if err := testCreateWorkflow(t, workflowID, tc.ProjectID, locationID); err != nil {
		t.Fatalf("testCreateWorkflow error: %v\n", err)
	}
	defer testCleanup(t, workflowID, tc.ProjectID, locationID)

	// Execute the workflow
	if err := executeWorkflowWithArguments(&buf, tc.ProjectID, workflowID, locationID); err != nil {
		t.Fatalf("executeWorkflow error: %v\n", err)
	}

	// Evaluate if the output contains the expected string.
	if got, want := buf.String(), "Execution results"; !strings.Contains(got, want) {
		t.Errorf("executeWorkflow: expected %q to contain %q", got, want)
	}

	// Evaluate if the output contains the argument value.
	if got, want := buf.String(), "Cloud"; !strings.Contains(got, want){
		t.Errorf("executeWorkflow: expected %q to contain %q", got, want)
	}
}

