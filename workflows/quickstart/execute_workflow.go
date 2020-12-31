// Copyright 2020 Google LLC
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
	"log"
	"os"
	"time"

	executions "cloud.google.com/go/workflows/executions/apiv1beta"
	executionspb "google.golang.org/genproto/googleapis/cloud/workflows/executions/v1beta"
)

// [START workflows_api_quickstart]

// executeWorkflow executes a workflow and returns the results from the workflow.
func executeWorkflow(projectID, locationID, workflowID string) (string, error) {
	ctx := context.Background()

	// Creates a client.
	client, err := executions.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return "", err
	}

	if workflowID == "" {
		workflowID = "myFirstWorkflow"
	}
	workflowPath := fmt.Sprintf("projects/%s/locations/%s/workflows/%s", projectID, locationID, workflowID)

	exe, err := client.CreateExecution(ctx, &executionspb.CreateExecutionRequest{
		Parent: workflowPath,
	})
	if err != nil {
		log.Fatalf("Failed to start workflow execution: %v", err)
		return "", err
	}

	name := exe.GetName()
	fmt.Fprintf(os.Stdout, "Created execution: %v\n", name)

	// Wait for execution to finish, then print results.
	finished := false
	backoffDelay := 1 * time.Second // Start wait with delay of 1s
	fmt.Println("Poll every second for result...")
	for !finished {
		e, err := client.GetExecution(ctx, &executionspb.GetExecutionRequest{
			Name: name,
		})
		if err != nil {
			log.Fatalf("Failed to get workflow execution: %v", err)
			return "", err
		}
		finished = e.State != executionspb.Execution_ACTIVE

		// If we haven't seen the result yet, wait a second.
		if !finished {
			fmt.Println("- Waiting for results...")
			time.Sleep(backoffDelay)
			backoffDelay *= 2 // Double the delay to provide exponential backoff.
		} else {
			fmt.Printf("Execution finished with state: %v\n", e.State)
			return e.Result, nil
		}
	}

	// should never happen
	return "", nil
}

// [END workflows_api_quickstart]
