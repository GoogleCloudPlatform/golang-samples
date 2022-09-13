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

// [START bigquery_create_table_partitioned]
import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
)

// createTablePartitioned demonstrates creating a table and specifying a time partitioning configuration.
func createTablePartitioned(projectID, datasetID, tableID string) error {
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
		{Name: "name", Type: bigquery.StringFieldType},
		{Name: "post_abbr", Type: bigquery.IntegerFieldType},
		{Name: "date", Type: bigquery.DateFieldType},
	}
	metadata := &bigquery.TableMetadata{
		TimePartitioning: &bigquery.TimePartitioning{
			Field:      "date",
			Expiration: 90 * 24 * time.Hour,
		},
		Schema: sampleSchema,
	}
	tableRef := client.Dataset(datasetID).Table(tableID)
	if err := tableRef.Create(ctx, metadata); err != nil {
		return err
	}
	return nil
}

// [END bigquery_create_table_partitioned]
