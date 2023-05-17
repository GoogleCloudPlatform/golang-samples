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

package snippets

// [START batch_get_task]
import (
	"context"
	"fmt"
	"io"

	batch "cloud.google.com/go/batch/apiv1"
	"cloud.google.com/go/batch/apiv1/batchpb"
)

// Retrieves the information about the specified job, most importantly its status
func getTask(w io.Writer, projectID, region, jobName, taskGroup string, taskNumber int32) error {
	// projectID := "your_project_id"
	// region := "us-central1"
	// jobName := "some-job"
	// taskGroup := "group0" // defaults to "group0" on job creation unless overridden
	// taskNumber := 0

	ctx := context.Background()
	batchClient, err := batch.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer batchClient.Close()

	req := &batchpb.GetTaskRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/jobs/%s/taskGroups/%s/tasks/%d",
			projectID, region, jobName, taskGroup, taskNumber),
	}

	response, err := batchClient.GetTask(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to get task: %w", err)
	}

	fmt.Fprintf(w, "Task info: %v\n", response)

	return nil
}

// [END batch_get_task]
