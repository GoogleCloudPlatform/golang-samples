// Copyright 2019 Google LLC
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

package howto

// [START job_search_delete_job]
import (
	"context"
	"fmt"
	"io"

	talent "cloud.google.com/go/talent/apiv4beta1"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
)

// deleteJob deletes an existing job by its resource name.
func deleteJob(w io.Writer, projectID, jobID string) error {
	ctx := context.Background()

	// Initialize a jobService client.
	c, err := talent.NewJobClient(ctx)
	if err != nil {
		return fmt.Errorf("talent.NewJobClient: %v", err)
	}

	// Construct a deleteJob request.
	jobName := fmt.Sprintf("projects/%s/jobs/%s", projectID, jobID)
	req := &talentpb.DeleteJobRequest{
		// The resource name of the job to be deleted.
		// The format is "projects/{project_id}/jobs/{job_id}".
		Name: jobName,
	}

	if err := c.DeleteJob(ctx, req); err != nil {
		return fmt.Errorf("Delete(%s): %v", jobName, err)
	}

	fmt.Fprintf(w, "Deleted job: %q\n", jobName)

	return err
}

// [END job_search_delete_job]
