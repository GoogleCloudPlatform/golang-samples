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

// [START bigquery_load_table_gcs_json]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// importJSONExplicitSchema demonstrates loading newline-delimited JSON data from Cloud Storage
// into a BigQuery table and providing an explicit schema for the data.
func importJSONExplicitSchema(projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}

	gcsRef := bigquery.NewGCSReference("gs://cloud-samples-data/bigquery/us-states/us-states.json")
	gcsRef.SourceFormat = bigquery.JSON
	gcsRef.Schema = bigquery.Schema{
		{Name: "name", Type: bigquery.StringFieldType},
		{Name: "post_abbr", Type: bigquery.StringFieldType},
	}
	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(gcsRef)
	loader.WriteDisposition = bigquery.WriteEmpty

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

// [END bigquery_load_table_gcs_json]
