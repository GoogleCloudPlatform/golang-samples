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

// [START bigquery_load_table_gcs_hivepartitioning]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// importWithHivePartitioning demonstrates loading data into a BigQuery managed
// table that includes data from a hive-based partition file layout.
func importWithHivePartitioning(projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	gcsRef := bigquery.NewGCSReference("gs://cloud-samples-data/bigquery/hive-partitioning-samples/customlayout/*")
	gcsRef.SourceFormat = bigquery.Parquet
	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(gcsRef)
	loader.HivePartitioningOptions = &bigquery.HivePartitioningOptions{
		Mode:            bigquery.CustomHivePartitioningMode,
		SourceURIPrefix: "gs://cloud-samples-data/bigquery/hive-partitioning-samples/customlayout/{pkey:STRING}/",
	}

	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}

	if status.Err() != nil {
		return fmt.Errorf("job completed with error: %w", status.Err())
	}
	return nil
}

// [END bigquery_load_table_gcs_hivepartitioning]
