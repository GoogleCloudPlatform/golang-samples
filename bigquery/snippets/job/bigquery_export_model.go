// Copyright 2022 Google LLC
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

package job

// [START bigquery_export_model]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// exportModel demonstrates how to export an existing
// BigQuery ML Model to Google Cloud Storage.
func exportModel(projectID, datasetID, modelID, gcsURI string) error {
	// projectID := "my-project-id"
	// datasetID := "dataset-id"
	// modelID := "model-id"
	// gcsURI := "gs://mybucket/path/to/model"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	gcsRef := bigquery.NewGCSReference(gcsURI)

	extractor := client.DatasetInProject(projectID, datasetID).Model(modelID).ExtractorTo(gcsRef)
	// You can choose to run the job in a specific location for more complex data locality scenarios.
	// Ex: In this example, source dataset and GCS bucket are in the US.
	extractor.Location = "US"

	job, err := extractor.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	return nil
}

// [END bigquery_export_model]
