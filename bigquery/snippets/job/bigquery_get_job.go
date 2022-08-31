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

// [START bigquery_get_job]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
)

// getJobInfo demonstrates retrieval of a job, which can be used to monitor
// completion or print metadata about the job.
func getJobInfo(w io.Writer, projectID, jobID string) error {
	// projectID := "my-project-id"
	// jobID := "my-job-id"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	job, err := client.JobFromID(ctx, jobID)
	if err != nil {
		return err
	}

	status := job.LastStatus()
	state := "Unknown"
	switch status.State {
	case bigquery.Pending:
		state = "Pending"
	case bigquery.Running:
		state = "Running"
	case bigquery.Done:
		state = "Done"
	}
	fmt.Fprintf(w, "Job %s was created %v and is in state %s\n",
		jobID, status.Statistics.CreationTime, state)
	return nil
}

// [END bigquery_get_job]
