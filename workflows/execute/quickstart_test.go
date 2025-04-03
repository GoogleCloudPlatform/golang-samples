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

package execute

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/workflows/v1"
)

func TestExecuteWorkflow(t *testing.T) {
	tc := testutil.SystemTest(t)

	workflowID := fmt.Sprintf("myFirstWorkflow_%s", uuid.NewString())
	locationID := "us-central1"

	var buf bytes.Buffer

	if err := testCreateWorkflow(t, workflowID, tc.ProjectID, locationID); err != nil {
		t.Fatalf("testCreateWorkflow error: %v\n", err)
	}
	defer testCleanup(t, workflowID, tc.ProjectID, locationID)

	if err := executeWorkflow(&buf, tc.ProjectID, workflowID, locationID); err != nil {
		t.Fatalf("executeWorkflow error: %v\n", err)
	}

	if got, want := buf.String(), "Execution results"; !strings.Contains(got, want) {
		t.Errorf("executeWorkflow: expected %q to contain %q", got, want)
	}

}

// testCreateWorkflow creates a testing workflow by the given name
func testCreateWorkflow(t *testing.T, workflowID, projectID, locationID string) error {
	t.Helper()

	ctx := context.Background()

	timeout := time.Minute * 5

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)
	workflowName := fmt.Sprintf("%s/workflows/%s", parent, workflowID)

	client, err := workflows.NewService(ctx, option.WithEndpoint("https://workflows.googleapis.com/"))
	if err != nil {
		return fmt.Errorf("workflows.NewService error: %w", err)
	}

	content, err := os.ReadFile("../myFirstWorkflow.yaml")
	if err != nil {
		return fmt.Errorf("os.ReadFile error: %w", err)
	}

	workflow := &workflows.Workflow{
		Name:           workflowName,
		SourceContents: string(content),
	}

	service := client.Projects.Locations.Workflows

	_, err = service.Create(parent, workflow).WorkflowId(workflowID).Do()
	if err != nil {
		return fmt.Errorf("service.Create failed to create workflow: %v", err)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

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

// testCreateWorkflow creates a testing workflow by the given name
func testCleanup(t *testing.T, workflowID, projectID, locationID string) error {
	t.Helper()

	ctx := context.Background()

	parent := fmt.Sprintf("projects/%s/locations/%s/workflows/%s", projectID, locationID, workflowID)

	client, err := workflows.NewService(ctx, option.WithEndpoint("https://workflows.googleapis.com/"))
	if err != nil {
		return fmt.Errorf("workflows.NewService error: %w", err)
	}

	service := client.Projects.Locations.Workflows
	_, err = service.Delete(parent).Do()
	if err != nil {
		return fmt.Errorf("service.Delete failed to delete workflow: %v", err)
	}

	return nil
}
