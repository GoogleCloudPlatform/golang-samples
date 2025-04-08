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

// [START bigquery_relax_column_query_append]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// relaxTableQuery demonstrates relaxing the schema of a table by appending query results to
// enable the table to allow NULL values.
func relaxTableQuery(projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	sampleSchema := bigquery.Schema{
		{Name: "full_name", Type: bigquery.StringFieldType, Required: true},
		{Name: "age", Type: bigquery.IntegerFieldType, Required: true},
	}
	meta := &bigquery.TableMetadata{
		Schema: sampleSchema,
	}
	tableRef := client.Dataset(datasetID).Table(tableID)
	if err := tableRef.Create(ctx, meta); err != nil {
		return err
	}
	// Now, append a query result that includes nulls, but allow the job to relax
	// all required columns.
	q := client.Query("SELECT \"Beyonce\" as full_name")
	q.QueryConfig.Dst = client.Dataset(datasetID).Table(tableID)
	q.SchemaUpdateOptions = []string{"ALLOW_FIELD_RELAXATION"}
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

// [END bigquery_relax_column_query_append]
