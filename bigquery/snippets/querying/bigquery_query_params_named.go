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

// [START bigquery_query_params_named]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// queryWithNamedParams demonstrate issuing a query using named query parameters.
func queryWithNamedParams(w io.Writer, projectID string) error {
	// projectID := "my-project-id"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	q := client.Query(
		`SELECT word, word_count
        FROM ` + "`bigquery-public-data.samples.shakespeare`" + `
        WHERE corpus = @corpus
        AND word_count >= @min_word_count
        ORDER BY word_count DESC;`)
	q.Parameters = []bigquery.QueryParameter{
		{
			Name:  "corpus",
			Value: "romeoandjuliet",
		},
		{
			Name:  "min_word_count",
			Value: 250,
		},
	}
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

// [END bigquery_query_params_named]
