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

// [START bigquery_query_stateless]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// queryStateless demonstrates issuing a query that may be run statelessly.
//
// To enable the stateless query preview feature, the QUERY_PREVIEW_ENABLED
// environmental variable should be set to `TRUE`.
func queryStateless(w io.Writer, projectID string) error {
	// projectID := "my-project-id"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// This example uses a nondeterministic query that doesn't reference
	// existing table(s).
	q := client.Query(
		"SELECT SESSION_USER() as whoami, CURRENT_TIMESTAMP as ts")
	// Run the query and process the returned row iterator.
	it, err := q.Read(ctx)
	if err != nil {
		return fmt.Errorf("query.Read(): %w", err)
	}

	// The iterator provides information about the source query execution.
	// Queries that were run in stateless mode will not have the source job
	// populated.
	if it.SourceJob() == nil {
		fmt.Fprintf(w, "Query was run statelessly.  Query ID: %q\n", it.QueryID())
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

// [END bigquery_query_stateless]
