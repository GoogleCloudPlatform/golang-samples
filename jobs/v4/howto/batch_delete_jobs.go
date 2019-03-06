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

// batchDeleteJobs deletes existing jobs by filter.
func batchDeleteJobs(w io.Writer, projectId string, filter string) error {
	ctx := context.Background()

	// Initialize a jobService client.
	c, err := talent.NewJobClient(ctx)
	if err != nil {
		fmt.Printf("talent.NewJobClient: %v", err)
		return err
	}

	// Construct a batchDeteleJobs request.
	req := &talentpb.BatchDeleteJobsRequest{
    Parent: "projects/" + projectId,
    // The fields eligible for filtering are `companyName` and `requisitionId`.
		Filter: filter,
	}

	if err := c.BatchDeleteJobs(ctx, req); err != nil {
		fmt.Printf("Batch deleting jobs from %s yielded: %v", filter, err)
		return err
	}

	fmt.Fprintf(w, "Batch deleted jobs from %s\n", filter)

	return err
}

// [END job_search_batch_delete_job]
