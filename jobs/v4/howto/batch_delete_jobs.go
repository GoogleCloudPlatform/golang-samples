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

import (
	"context"
	"fmt"
	"io"

	talent "cloud.google.com/go/talent/apiv4beta1"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
)

// [START job_search_batch_delete_job]

// deleteJob deletes an existing job by filter.
func batchDeleteJobs(w io.Writer, projectID string, filter string) error {
	ctx := context.Background()

	// Create a job service client.
	c, err := talent.NewJobClient(ctx)
	if err != nil {
		return fmt.Errorf("talent.NewJobClient: %v", err)
	}

	// Construct a GetJobRequest.
	req := &talentpb.BatchDeleteJobsRequest{
    Parent: "projects/" + projectID,
    // The fields eligible for filtering are `companyName` and `requisitionId`.
		Filter: filter,
	}

	if err := c.BatchDeleteJobs(ctx, req); err != nil {
		return fmt.Errorf("Delete jobs from %s: %v", filter, err)
	}

	return err
}

// [END job_search_batch_delete_job]
