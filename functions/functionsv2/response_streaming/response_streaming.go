// Copyright 2023 Google LLC
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

// [START functions_response_streaming]

// Package responsestreaming contains a function that streams out a large payload in chunks.
package responsestreaming

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"google.golang.org/api/iterator"
)

func init() {
	functions.HTTP("streamBigQuery", streamBigQuery)
}

// streamBigQuery retrieves a large payload from BigQuery public dataset and streams its rows.
func streamBigQuery(w http.ResponseWriter, r *http.Request) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		fmt.Println("GOOGLE_CLOUD_PROJECT environment variable must be set.")
		os.Exit(1)
	}
	ctx := context.Background()
	// Must include project ID in client.
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	// Run the query and ensure the job finishes without error.
	rows, err := query(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	// Stream out the payload by iterating rows and flushing out the writer.
	streamResults(w, rows)
}

// query returns a row iterator suitable for reading query results.
func query(ctx context.Context, client *bigquery.Client) (*bigquery.RowIterator, error) {
	q := client.Query(
		"SELECT abstract FROM `bigquery-public-data.breathe.bioasq` LIMIT 1000")
	q.Location = "US"
	return q.Read(ctx)
}

// streamResults streams results from a query of BigQuery's public dataset.
func streamResults(w io.Writer, it *bigquery.RowIterator) {
	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		fmt.Fprintln(w, row[0])
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
			fmt.Fprintln(w, "Successfully flushed row")
		}
	}
}

// [END functions_response_streaming]
