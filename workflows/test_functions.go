// Copyright 2025 Google LLC
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

package workflows

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/workflows/v1"
)

// testCreateWorkflow creates a testing workflow by the given name.
func testCreateWorkflow(t *testing.T, workflowID, projectID, locationID string) error {
	t.Helper()

	ctx := context.Background()

	timeout := time.Minute * 5

	// Build the parent and workflowName by the given values.
	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)
	workflowName := fmt.Sprintf("%s/workflows/%s", parent, workflowID)

	// Create client
	client, err := workflows.NewService(ctx)
	if err != nil {
		return fmt.Errorf("workflows.NewService error: %w", err)
	}

	// Read file's content
	content, err := os.ReadFile("myFirstWorkflow.yaml")
	if err != nil {
		return fmt.Errorf("os.ReadFile error: %w", err)
	}

	// Build the workflow by assigning the name and the content.
	workflow := &workflows.Workflow{
		Name:           workflowName,
		SourceContents: string(content),
	}

	// Get the workflow service for avoiding repetitive code.
	service := client.Projects.Locations.Workflows

	// Create workflow
	_, err = service.Create(parent, workflow).WorkflowId(workflowID).Do()
	if err != nil {
		return fmt.Errorf("service.Create failed to create workflow: %v", err)
	}

	// To avoid an infinite loop a context with timeout will be created.
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Loop the workflow's creation state until it's different from active.
	for workflow.State != "ACTIVE" {
		workflow, err = service.Get(workflowName).Do()
		if err != nil {
			return fmt.Errorf("service.Get.Do error: %w", err)
		}

		select {
		case <-time.After(time.Second * 1):
		case <-timeoutCtx.Done():
			return timeoutCtx.Err()
		}
	}

	return nil
}

// testCleanup deletes a workflow to ensure that all resources
// created or utilized during the execution of the test are
// released.
func testCleanup(t *testing.T, workflowID, projectID, locationID string) error {
	t.Helper()

	ctx := context.Background()

	// Build the parent and workflowName by the given values.
	parent := fmt.Sprintf("projects/%s/locations/%s/workflows/%s", projectID, locationID, workflowID)

	// Create client.
	client, err := workflows.NewService(ctx)
	if err != nil {
		return fmt.Errorf("workflows.NewService error: %w", err)
	}

	// Delete workflow.
	service := client.Projects.Locations.Workflows
	_, err = service.Delete(parent).Do()
	if err != nil {
		return fmt.Errorf("service.Delete failed to delete workflow: %v", err)
	}

	return nil
}

// testGenerateWorkflowID creates a random string to
// assign it to a new workflow.
func testGenerateWorkflowID() string {
	return fmt.Sprintf("workflow-go-test%s", uuid.NewString())
}
