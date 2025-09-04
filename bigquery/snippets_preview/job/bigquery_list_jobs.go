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

package job

// [START bigquery_list_jobs_preview]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery/v2/apiv2/bigquerypb"
	"cloud.google.com/go/bigquery/v2/apiv2_client"

	"google.golang.org/api/iterator"
)

// listJobs demonstrates iterating through job metadata.
func listJobs(client *apiv2_client.Client, w io.Writer, projectID string) error {
	// client can be instantiated per-RPC service, or use cloud.google.com/go/bigquery/v2/apiv2_client to create
	// an aggregate client.
	//
	// projectID := "my-project-id"
	ctx := context.Background()

	req := &bigquerypb.ListJobsRequest{
		ProjectId: projectID,
		// Only list pending or running jobs.
		StateFilter: []bigquerypb.ListJobsRequest_StateFilter{
			bigquerypb.ListJobsRequest_PENDING,
			bigquerypb.ListJobsRequest_RUNNING,
		},
	}

	// ListJobs returns an iterator so users don't have to manage pagination when processing
	// the results.
	it := client.ListJobs(ctx, req)

	// Process data from the iterator one result at a time, and stop after we
	// process a fixed number of jobs.  While the number of inflight
	// (pending or running) jobs may be more reasonable, listing all jobs can yield
	// a potentially very large number of results.
	maxJobs := 10
	for numJobs := 0; numJobs < maxJobs; numJobs++ {
		job, err := it.Next()
		if err == iterator.Done {
			// We're reached the end of the iteration, break the loop.
			break
		}
		if err != nil {
			return fmt.Errorf("iterator errored: %w", err)
		}
		// Print basic information to the provided writer.
		fmt.Fprintf(w, "job %q in location %q is in state %q\n",
			job.GetJobReference().GetJobId(),
			job.GetJobReference().GetLocation(),
			job.GetState())
	}
	return nil
}

// [END bigquery_list_jobs_preview]
