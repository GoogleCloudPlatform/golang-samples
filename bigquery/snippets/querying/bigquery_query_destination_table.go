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

// [START bigquery_query_destination_table]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// queryWithDestination demonstrates saving the results of a query to a specific table by setting the destination
// via the API properties.
func queryWithDestination(w io.Writer, projectID, destDatasetID, destTableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	q := client.Query("SELECT 17 as my_col")
	q.Location = "US" // Location must match the dataset(s) referenced in query.
	q.QueryConfig.Dst = client.Dataset(destDatasetID).Table(destTableID)
	// Run the query and process the returned row iterator.
	it, err := q.Read(ctx)
	if err != nil {
		return fmt.Errorf("query.Read(): %w", err)
	}
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

// [END bigquery_query_destination_table]
