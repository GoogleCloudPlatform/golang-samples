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

package dataset

// [START bigquery_get_dataset_labels]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
)

// printDatasetLabels retrieves label metadata from a dataset and prints it to an io.Writer.
func printDatasetLabels(w io.Writer, projectID, datasetID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	meta, err := client.Dataset(datasetID).Metadata(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Dataset %s labels:\n", datasetID)
	if len(meta.Labels) == 0 {
		fmt.Fprintln(w, "Dataset has no labels defined.")
		return nil
	}
	for k, v := range meta.Labels {
		fmt.Fprintf(w, "\t%s:%s\n", k, v)
	}
	return nil
}

// [END bigquery_get_dataset_labels]
