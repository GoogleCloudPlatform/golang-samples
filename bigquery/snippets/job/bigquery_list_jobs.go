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

// [START bigquery_list_jobs]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// listJobs demonstrates iterating through the BigQuery jobs collection.
func listJobs(w io.Writer, projectID string) error {
	// projectID := "my-project-id"
	// jobID := "my-job-id"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	it := client.Jobs(ctx)
	// List up to 10 jobs to demonstrate iteration.
	for i := 0; i < 10; i++ {
		j, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		state := "Unknown"
		switch j.LastStatus().State {
		case bigquery.Pending:
			state = "Pending"
		case bigquery.Running:
			state = "Running"
		case bigquery.Done:
			state = "Done"
		}
		fmt.Fprintf(w, "Job %s in state %s\n", j.ID(), state)
	}
	return nil
}

// [END bigquery_list_jobs]
