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

// [START bigquery_get_table_labels]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
)

// tableLabels demonstrates fetching metadata from a table and printing the Label metadata to an io.Writer.
func tableLabels(w io.Writer, projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	meta, err := client.Dataset(datasetID).Table(tableID).Metadata(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Table %s labels:\n", datasetID)
	if len(meta.Labels) == 0 {
		fmt.Fprintln(w, "Table has no labels defined.")
		return nil
	}
	for k, v := range meta.Labels {
		fmt.Fprintf(w, "\t%s:%s\n", k, v)
	}
	return nil
}

// [END bigquery_get_table_labels]
