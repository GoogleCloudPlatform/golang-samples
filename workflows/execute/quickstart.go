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

package execute

// [START workflows_api_quickstart]

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	workflowexecutions "google.golang.org/api/workflowexecutions/v1"
)

// Execute a workflow and print the execution results.
//
// A workflow consists of a series of steps described
// using the Workflows syntax, and can be written in either YAML or JSON.
//
// For more information about Workflows see:
// https://cloud.google.com/workflows/docs/overview
func executeWorkflow(w io.Writer, projectID, workflowID, locationID string) error {
	// TODO(developer): uncomment and update the following lines:
	// projectID := "my-project-id"
	// workflowID := "my-workflow-id"
	// locationID := "location-id"

	ctx := context.Background()

	delay := time.Second * 1
	timeout := time.Minute * 10 // Timeout for workflow execution

	// Construct the location path.
	parent := fmt.Sprintf("projects/%s/locations/%s/workflows/%s", projectID, locationID, workflowID)

	// Create execution client.
	client, err := workflowexecutions.NewService(ctx)
	if err != nil {
		return fmt.Errorf("workflowexecutions.NewService error: %w", err)
	}

	// Get execution service.
	service := client.Projects.Locations.Workflows.Executions

	// Build and run the new workflow execution.
	res, err := service.Create(parent, &workflowexecutions.Execution{}).Do()
	if err != nil {
		return fmt.Errorf("service.Create.Do error: %w", err)
	}
	log.Printf("- Execution started...")

	// Set up a context with timeout to prevent an infinite loop.
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for res.State == "ACTIVE" {
		// Request for getting the updated state of the execution.
		getReq := service.Get(res.Name)
		res, err = getReq.Do()
		if err != nil {
			return fmt.Errorf("getReq error: %w", err)
		}

		// Double the delay to provide exponential backoff.
		delay *= 2

		// Select statement to Evaluate whether the delay or the context with timeout has completed.
		select {
		case <-time.After(delay):
			log.Println("- Waiting for results...")
		case <-timeoutCtx.Done():
			return timeoutCtx.Err()
		}
	}

	fmt.Fprintf(w, "Execution finished with state: %s\n", res.State)
	fmt.Fprintf(w, "Execution results: %s\n", res.Result)

	return nil
}

// [END workflows_api_quickstart]
