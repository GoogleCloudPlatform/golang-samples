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

package table

// [START bigquery_get_table]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
)

// printTableInfo demonstrates fetching metadata from a table and printing some basic information
// to an io.Writer.
func printTableInfo(w io.Writer, projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	meta, err := client.Dataset(datasetID).Table(tableID).Metadata(ctx)
	if err != nil {
		return err
	}
	// Print basic information about the table.
	fmt.Fprintf(w, "Schema has %d top-level fields\n", len(meta.Schema))
	fmt.Fprintf(w, "Description: %s\n", meta.Description)
	fmt.Fprintf(w, "Rows in managed storage: %d\n", meta.NumRows)
	return nil
}

// [END bigquery_get_table]
