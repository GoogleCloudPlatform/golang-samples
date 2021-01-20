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
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"

	workflows "cloud.google.com/go/workflows/apiv1beta"
	workflowspb "google.golang.org/genproto/googleapis/cloud/workflows/v1beta"
)

func TestExecuteWorkflow(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := "us-central1"
	workflowName := "myFirstWorkflow"

	err := deployWorkflow(tc.ProjectID, locationID, workflowName)
	if err != nil {
		t.Errorf("ExecuteWorkflow(deploy): %v", err)
	}

	_, err = executeWorkflow(tc.ProjectID, locationID, workflowName)
	if err != nil {
		t.Errorf("ExecuteWorkflow(execute): %v", err)
	}
}

// deployWorkflow deploys a workflow.
func deployWorkflow(projectID, locationID, workflowID string) error {
	workflowExists, err := workflowExists(projectID, locationID, workflowID)
	if workflowExists == true && err != nil {
		return fmt.Errorf("deployWorkflow: %v", err)
	}

	ctx := context.Background()
	client, err := workflows.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("workflows.NewClient: %v", err)
	}
	workflowPath := fmt.Sprintf("projects/%s/locations/%s/workflows/%s", projectID, locationID, workflowID)
	print(workflowPath)
	_, err = client.CreateWorkflow(ctx, &workflowspb.CreateWorkflowRequest{
		Parent: workflowPath,
		Workflow: &workflowspb.Workflow{
			Name: workflowID,
			// Copied from:
			// https://github.com/GoogleCloudPlatform/workflows-samples/blob/main/src/myFirstWorkflow.workflows.yaml
			SourceCode: &workflowspb.Workflow_SourceContents{
				SourceContents: "# Copyright 2020 Google LLC\n#\n# Licensed under the Apache License, Version 2.0 (the \"License\");\n# you may not use this file except in compliance with the License.\n# You may obtain a copy of the License at\n#\n#      http://www.apache.org/licenses/LICENSE-2.0\n#\n# Unless required by applicable law or agreed to in writing, software\n# distributed under the License is distributed on an \"AS IS\" BASIS,\n# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n# See the License for the specific language governing permissions and\n# limitations under the License.\n\n# [START workflows_myfirstworkflow]\n- getCurrentTime:\n    call: http.get\n    args:\n      url: https://us-central1-workflowsample.cloudfunctions.net/datetime\n    result: currentTime\n- readWikipedia:\n    call: http.get\n    args:\n      url: https://en.wikipedia.org/w/api.php\n      query:\n        action: opensearch\n        search: ${currentTime.body.dayOfTheWeek}\n    result: wikiResult\n- returnResult:\n    return: ${wikiResult.body[1]}\n# [END workflows_myfirstworkflow]\n",
			},
		},
		WorkflowId: workflowID,
	})
	return fmt.Errorf("client.CreateWorkflow: %v", err)
}

func workflowExists(projectID, locationID, workflowID string) (bool, error) {
	ctx := context.Background()
	client, err := workflows.NewClient(ctx)
	if err != nil {
		return false, fmt.Errorf("workflows.NewClient: %v", err)
	}
	workflowPath := fmt.Sprintf("projects/%s/locations/%s/workflows/%s", projectID, locationID, workflowID)
	wf, err := client.GetWorkflow(ctx, &workflowspb.GetWorkflowRequest{
		Name: workflowPath,
	})
	if err != nil {
		return false, fmt.Errorf("client.GetWorkflow: %v", err)
	}
	return wf.State == workflowspb.Workflow_ACTIVE, nil
}
