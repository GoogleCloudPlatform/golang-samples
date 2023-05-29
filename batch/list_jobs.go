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

// [START batch_list_jobs]
import (
	"context"
	"fmt"
	"io"

	batch "cloud.google.com/go/batch/apiv1"
	"cloud.google.com/go/batch/apiv1/batchpb"
	"google.golang.org/api/iterator"
)

// Lists all jobs in the given project and region
func listJobs(w io.Writer, projectID, region string) error {
	// projectID := "your_project_id"
	// region := "us-central1"

	ctx := context.Background()
	batchClient, err := batch.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer batchClient.Close()

	req := &batchpb.ListJobsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, region),
	}

	var jobs []*batchpb.Job
	it := batchClient.ListJobs(ctx, req)

	for {
		job, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("unable to list jobs: %w", err)
		}
		jobs = append(jobs, job)
	}

	fmt.Fprintf(w, "Jobs: %v\n", jobs)

	return nil
}

// [END batch_list_jobs]
