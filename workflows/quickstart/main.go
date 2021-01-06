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

// Simple CLI to run the executeWorkflow function.
// Used for one-off testing and development.

// [START workflows_api_quickstart]

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	executions "cloud.google.com/go/workflows/executions/apiv1beta"
	executionspb "google.golang.org/genproto/googleapis/cloud/workflows/executions/v1beta"
)

// executeWorkflow executes a workflow and returns the results from the workflow.
func executeWorkflow(projectID, locationID, workflowID string) (string, error) {
	ctx := context.Background()

	// Creates a client.
	client, err := executions.NewClient(ctx)
	if err != nil {
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
		return "", err
	}

	name := exe.GetName()
	fmt.Fprintf(os.Stdout, "Created execution: %v\n", name)

	// Wait for execution to finish, then print results.
	backoffDelay := 1 * time.Second // Start wait with delay of 1s.
	fmt.Println("Poll for result...")
	for {
		e, err := client.GetExecution(ctx, &executionspb.GetExecutionRequest{
			Name: name,
		})
		if err != nil {
			return "", err
		}

		// If we haven't seen the result yet, wait a second.
		if e.State == executionspb.Execution_ACTIVE {
			fmt.Printf("- Waiting %ds for results...\n", backoffDelay/time.Second)
			time.Sleep(backoffDelay)
			backoffDelay *= 2 // Double the delay to provide exponential backoff.
		} else {
			fmt.Printf("Execution finished with state: %v\n", e.State)
			return e.Result, nil
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: Must include 2 arguments for projectID, locationID")
		os.Exit(1)
	}
	projectID := os.Args[1]
	locationID := os.Args[2]
	workflowID := ""
	if len(os.Args) > 3 {
		workflowID = os.Args[3]
	}

	res, err := executeWorkflow(projectID, locationID, workflowID)
	if err != nil {
		log.Fatalf("Failure in workflow execution: %v", err)
	}
	var jsonStringArr []string
	err = json.Unmarshal([]byte(res), &jsonStringArr)

	fmt.Print(jsonStringArr)
}

// [END workflows_api_quickstart]
