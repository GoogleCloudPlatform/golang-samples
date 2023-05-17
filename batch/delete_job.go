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

// [START batch_delete_job]
import (
	"context"
	"fmt"
	"io"

	batch "cloud.google.com/go/batch/apiv1"
	"cloud.google.com/go/batch/apiv1/batchpb"
)

// Deletes the specified job
func deleteJob(w io.Writer, projectID, region, jobName string) error {
	// projectID := "your_project_id"
	// region := "us-central1"
	// jobName := "some-job"

	ctx := context.Background()
	batchClient, err := batch.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer batchClient.Close()

	req := &batchpb.DeleteJobRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectID, region, jobName),
	}

	response, err := batchClient.DeleteJob(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to delete job: %w", err)
	}

	fmt.Fprintf(w, "Job deleted: %v\n", response)

	return nil
}

// [END batch_delete_job]
