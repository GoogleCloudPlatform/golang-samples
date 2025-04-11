// Copyright 2020 Google LLC
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

// [START bigquery_create_table_external_hivepartitioned]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// createTableExternalHivePartitioned demonstrates creating an external table with hive partitioning.
func createTableExternalHivePartitioned(projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydatasetid"
	// tableID := "mytableid"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// First, we'll define table metadata to represent a table that's backed by parquet files held in
	// Cloud Storage.
	//
	// Example file:
	// gs://cloud-samples-data/bigquery/hive-partitioning-samples/autolayout/dt=2020-11-15/file1.parquet
	metadata := &bigquery.TableMetadata{
		Description: "An example table that demonstrates hive partitioning against external parquet files",
		ExternalDataConfig: &bigquery.ExternalDataConfig{
			SourceFormat: bigquery.Parquet,
			SourceURIs:   []string{"gs://cloud-samples-data/bigquery/hive-partitioning-samples/autolayout/*"},
			AutoDetect:   true,
		},
	}

	// The layout of the files in here is compatible with the layout requirements for hive partitioning,
	// so we can add an optional Hive partitioning configuration to leverage the object paths for deriving
	// partitioning column information.
	//
	// For more information on how partitions are extracted, see:
	// https://cloud.google.com/bigquery/docs/hive-partitioned-queries-gcs
	//
	// We have a "/dt=YYYY-MM-DD/" path component in our example files as documented above.  Autolayout will
	// expose this as a column named "dt" of type DATE.
	metadata.ExternalDataConfig.HivePartitioningOptions = &bigquery.HivePartitioningOptions{
		Mode:                   bigquery.AutoHivePartitioningMode,
		SourceURIPrefix:        "gs://cloud-samples-data/bigquery/hive-partitioning-samples/autolayout/",
		RequirePartitionFilter: true,
	}

	// Create the external table.
	tableRef := client.Dataset(datasetID).Table(tableID)
	if err := tableRef.Create(ctx, metadata); err != nil {
		return fmt.Errorf("table creation failure: %w", err)
	}
	return nil
}

// [END bigquery_create_table_external_hivepartitioned]
