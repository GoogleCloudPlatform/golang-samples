// Copyright 2024 Google LLC
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

// [START bigquery_query_job_optional]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// queryJobOptional demonstrates issuing a query that doesn't require a
// a corresponding job.
func queryJobOptional(w io.Writer, projectID string) error {
	// projectID := "my-project-id"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	q := client.Query(`
		SELECT
  			name, gender,
  			SUM(number) AS total
		FROM
			bigquery-public-data.usa_names.usa_1910_2013
		GROUP BY 
			name, gender
		ORDER BY
			total DESC
		LIMIT 10
		`)
	// Run the query and process the returned row iterator.
	it, err := q.Read(ctx)
	if err != nil {
		return fmt.Errorf("query.Read(): %w", err)
	}

	// The iterator provides information about the query execution.
	// Queries that were run in short query mode will not have the source job
	// populated.
	if it.SourceJob() == nil {
		fmt.Fprintf(w, "Query was run in short mode.  Query ID: %q\n", it.QueryID())
	} else {
		j := it.SourceJob()
		qualifiedJobID := fmt.Sprintf("%s:%s.%s", j.ProjectID(), j.Location(), j.ID())
		fmt.Fprintf(w, "Query was run with job state.  Job ID: %q, Query ID: %q\n",
			qualifiedJobID, it.QueryID())
	}

	// Print row data.
	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(w, row)
	}
	return nil
}

// [END bigquery_query_job_optional]
