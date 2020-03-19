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

package loadingdata

// [START bigquery_load_table_gcs_csv_truncate]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// importCSVTruncate demonstrates loading data from CSV data in Cloud Storage and overwriting/truncating
// data in the existing table.
func importCSVTruncate(projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	gcsRef := bigquery.NewGCSReference("gs://cloud-samples-data/bigquery/us-states/us-states.csv")
	gcsRef.SourceFormat = bigquery.CSV
	gcsRef.AutoDetect = true
	gcsRef.SkipLeadingRows = 1
	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(gcsRef)
	loader.WriteDisposition = bigquery.WriteTruncate

	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}

	if status.Err() != nil {
		return fmt.Errorf("job completed with error: %v", status.Err())
	}
	return nil
}

// [END bigquery_load_table_gcs_csv_truncate]
