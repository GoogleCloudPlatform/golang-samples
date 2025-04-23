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
	"context"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// TestExecuteWorkflowWithArguments creates a workflow, executes it, and
// validates whether the execution results are returned correctly.
func TestExecuteWorkflowWithArguments(t *testing.T) {
	tc := testutil.SystemTest(t)

	workflowID := generateWorkflowID()
	locationID := "us-central1"

	var err error
	var buf bytes.Buffer

	// Create the test workflow that will be cleaned up once the test is done.
	if err = createWorkflow(t, workflowID, tc.ProjectID, locationID); err != nil {
		t.Fatalf("testCreateWorkflow error: %v\n", err)
	}
	defer cleanup(t, workflowID, tc.ProjectID, locationID)

	// Execute the workflow with a timeout of 10 minutes.
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	chanErr := make(chan error, 1) // Buffered channel for receiving the function's result

	// Goroutine that expects the returning value from the workflow execution and sends it to the channel.
	go func() {
		chanErr <- executeWorkflowWithArguments(&buf, tc.ProjectID, workflowID, locationID)
		close(chanErr)
	}()

	// Block until timeout or a return value is received from the function call.
	select {
	case <-ctxTimeout.Done():
		close(chanErr)
		t.Fatalf("executeWorkflow error: %v", context.DeadlineExceeded)
	case err = <-chanErr:
		if err != nil {
			t.Fatalf("executeWorkflow error: %v\n", err)
		}
	}

	// Evaluate whether the output contains the expected string.
	if got, want := buf.String(), "Execution results"; !strings.Contains(got, want) {
		t.Errorf("executeWorkflow: expected %q to contain %q", got, want)
	}

	// Evaluate if the output contains the argument value.
	if got, want := buf.String(), "Cloud"; !strings.Contains(got, want) {
		t.Errorf("executeWorkflow: expected %q to contain %q", got, want)
	}
}
