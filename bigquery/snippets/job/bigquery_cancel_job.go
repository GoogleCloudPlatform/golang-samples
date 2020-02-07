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

package job

// [START bigquery_cancel_job]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// cancelJob demonstrates how a job cancellation request can be issued for a specific
// BigQuery job.
func cancelJob(projectID, jobID string) error {
	// projectID := "my-project-id"
	// jobID := "my-job-id"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	job, err := client.JobFromID(ctx, jobID)
	if err != nil {
		return nil
	}
	return job.Cancel(ctx)
}

// [END bigquery_cancel_job]
