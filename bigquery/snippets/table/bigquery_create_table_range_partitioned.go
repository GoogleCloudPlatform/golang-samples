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

package table

// [START bigquery_create_table_range_partitioned]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// createTableRangeParitioned demonstrates creating a table and specifying a
// range partitioning configuration.
func createTableRangePartitioned(projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydatasetid"
	// tableID := "mytableid"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	sampleSchema := bigquery.Schema{
		{Name: "full_name", Type: bigquery.StringFieldType},
		{Name: "city", Type: bigquery.StringFieldType},
		{Name: "zipcode", Type: bigquery.IntegerFieldType},
	}

	metadata := &bigquery.TableMetadata{
		RangePartitioning: &bigquery.RangePartitioning{
			Field: "zipcode",
			Range: &bigquery.RangePartitioningRange{
				Start:    0,
				End:      100000,
				Interval: 10,
			},
		},
		Schema: sampleSchema,
	}
	tableRef := client.Dataset(datasetID).Table(tableID)
	if err := tableRef.Create(ctx, metadata); err != nil {
		return err
	}
	return nil
}

// [END bigquery_create_table_range_partitioned]
