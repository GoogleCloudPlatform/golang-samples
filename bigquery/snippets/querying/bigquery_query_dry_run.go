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

// [START bigquery_query_dry_run]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
)

// queryDryRun demonstrates issuing a dry run query to validate query structure and
// provide an estimate of the bytes scanned.
func queryDryRun(w io.Writer, projectID string) error {
	// projectID := "my-project-id"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	q := client.Query(`
	SELECT
		name,
		COUNT(*) as name_count
	FROM ` + "`bigquery-public-data.usa_names.usa_1910_2013`" + `
	WHERE state = 'WA'
	GROUP BY name`)
	q.DryRun = true
	// Location must match that of the dataset(s) referenced in the query.
	q.Location = "US"

	job, err := q.Run(ctx)
	if err != nil {
		return err
	}
	// Dry run is not asynchronous, so get the latest status and statistics.
	status := job.LastStatus()
	if err := status.Err(); err != nil {
		return err
	}
	fmt.Fprintf(w, "This query will process %d bytes\n", status.Statistics.TotalBytesProcessed)
	return nil
}

// [END bigquery_query_dry_run]
