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

package models

// [START bigquery_get_model]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
)

// printModelInfo demonstrates fetching metadata about a BigQuery ML model and printing some of
// it to an io.Writer.
func printModelInfo(w io.Writer, projectID, datasetID, modelID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// modelID := "mymodel"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	meta, err := client.Dataset(datasetID).Model(modelID).Metadata(ctx)
	if err != nil {
		return fmt.Errorf("couldn't retrieve metadata: %w", err)
	}
	fmt.Fprintf(w, "Got model '%q' with friendly name '%q'\n", modelID, meta.Name)
	return nil
}

// [END bigquery_get_model]
