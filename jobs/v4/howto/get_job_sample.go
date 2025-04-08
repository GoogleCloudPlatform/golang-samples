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

// [START job_search_get_job]
import (
	"context"
	"fmt"
	"io"

	talent "cloud.google.com/go/talent/apiv4beta1"
	"cloud.google.com/go/talent/apiv4beta1/talentpb"
)

// getJob gets an existing job by its resource name.
func getJob(w io.Writer, projectID, jobID string) (*talentpb.Job, error) {
	ctx := context.Background()

	// Initialize a jobService client.
	c, err := talent.NewJobClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("talent.NewJobClient: %w", err)
	}
	defer c.Close()

	// Construct a getJob request.
	jobName := fmt.Sprintf("projects/%s/jobs/%s", projectID, jobID)
	req := &talentpb.GetJobRequest{
		// The resource name of the job to retrieve.
		// The format is "projects/{project_id}/jobs/{job_id}".
		Name: jobName,
	}

	resp, err := c.GetJob(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetJob: %w", err)
	}

	fmt.Fprintf(w, "Job: %q\n", resp.GetName())
	fmt.Fprintf(w, "Job title: %v\n", resp.GetTitle())

	return resp, err
}

// [END job_search_get_job]
