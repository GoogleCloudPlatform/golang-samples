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

// [START bigquery_create_table_clustered]
import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
)

// createTableClustered demonstrates creating a BigQuery table with advanced properties like
// partitioning and clustering features.
func createTableClustered(projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydatasetid"
	// tableID := "mytableid"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	sampleSchema := bigquery.Schema{
		{Name: "timestamp", Type: bigquery.TimestampFieldType},
		{Name: "origin", Type: bigquery.StringFieldType},
		{Name: "destination", Type: bigquery.StringFieldType},
		{Name: "amount", Type: bigquery.NumericFieldType},
	}
	metaData := &bigquery.TableMetadata{
		Schema: sampleSchema,
		TimePartitioning: &bigquery.TimePartitioning{
			Field:      "timestamp",
			Expiration: 90 * 24 * time.Hour,
		},
		Clustering: &bigquery.Clustering{
			Fields: []string{"origin", "destination"},
		},
	}
	tableRef := client.Dataset(datasetID).Table(tableID)
	if err := tableRef.Create(ctx, metaData); err != nil {
		return err
	}
	return nil
}

// [END bigquery_create_table_clustered]
