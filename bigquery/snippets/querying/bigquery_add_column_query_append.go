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

package querying

// [START bigquery_add_column_query_append]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// createTableAndWidenQuery demonstrates how the schema of a table can be modified to add columns by appending
// query results that include the new columns.
func createTableAndWidenQuery(projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// First, we create a sample table.
	sampleSchema := bigquery.Schema{
		{Name: "full_name", Type: bigquery.StringFieldType, Required: true},
		{Name: "age", Type: bigquery.IntegerFieldType, Required: true},
	}
	original := &bigquery.TableMetadata{
		Schema: sampleSchema,
	}
	tableRef := client.Dataset(datasetID).Table(tableID)
	if err := tableRef.Create(ctx, original); err != nil {
		return err
	}
	// Our table has two columns.  We'll introduce a new favorite_color column via
	// a subsequent query that appends to the table.
	q := client.Query("SELECT \"Timmy\" as full_name, 85 as age, \"Blue\" as favorite_color")
	q.SchemaUpdateOptions = []string{"ALLOW_FIELD_ADDITION"}
	q.QueryConfig.Dst = client.Dataset(datasetID).Table(tableID)
	q.WriteDisposition = bigquery.WriteAppend
	q.Location = "US"
	job, err := q.Run(ctx)
	if err != nil {
		return err
	}
	_, err = job.Wait(ctx)
	if err != nil {
		return err
	}
	return nil
}

// [END bigquery_add_column_query_append]
