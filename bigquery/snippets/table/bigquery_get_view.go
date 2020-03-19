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

// [START bigquery_get_view]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
)

// getView demonstrates fetching the metadata from a BigQuery logical view and printing it to an io.Writer.
func getView(w io.Writer, projectID, datasetID, viewID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// viewID := "myview"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	view := client.Dataset(datasetID).Table(viewID)
	meta, err := view.Metadata(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "View %s, query: %s\n", view.FullyQualifiedName(), meta.ViewQuery)
	return nil
}

// [END bigquery_get_view]
