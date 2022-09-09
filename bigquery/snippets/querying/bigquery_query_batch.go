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

package querying

// [START bigquery_query_batch]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/bigquery"
)

// queryBatch demonstrates issuing a query job using batch priority.
func queryBatch(w io.Writer, projectID, dstDatasetID, dstTableID string) error {
	// projectID := "my-project-id"
	// dstDatasetID := "mydataset"
	// dstTableID := "mytable"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// Build an aggregate table.
	q := client.Query(`
		SELECT
  			corpus,
  			SUM(word_count) as total_words,
  			COUNT(1) as unique_words
		FROM ` + "`bigquery-public-data.samples.shakespeare`" + `
		GROUP BY corpus;`)
	q.Priority = bigquery.BatchPriority
	q.QueryConfig.Dst = client.Dataset(dstDatasetID).Table(dstTableID)

	// Start the job.
	job, err := q.Run(ctx)
	if err != nil {
		return err
	}
	// Job is started and will progress without interaction.
	// To simulate other work being done, sleep a few seconds.
	time.Sleep(5 * time.Second)
	status, err := job.Status(ctx)
	if err != nil {
		return err
	}

	state := "Unknown"
	switch status.State {
	case bigquery.Pending:
		state = "Pending"
	case bigquery.Running:
		state = "Running"
	case bigquery.Done:
		state = "Done"
	}
	// You can continue to monitor job progress until it reaches
	// the Done state by polling periodically.  In this example,
	// we print the latest status.
	fmt.Fprintf(w, "Job %s in Location %s currently in state: %s\n", job.ID(), job.Location(), state)

	return nil

}

// [END bigquery_query_batch]
