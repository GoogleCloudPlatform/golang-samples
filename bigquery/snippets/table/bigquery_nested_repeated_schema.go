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

// [START bigquery_nested_repeated_schema]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
)

// createTableComplexSchema demonstrates creating a BigQuery table and specifying a complex schema that includes
// an array of Struct types.
func createTableComplexSchema(w io.Writer, projectID, datasetID, tableID string) error {
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
		{Name: "id", Type: bigquery.StringFieldType},
		{Name: "first_name", Type: bigquery.StringFieldType},
		{Name: "last_name", Type: bigquery.StringFieldType},
		{Name: "dob", Type: bigquery.DateFieldType},
		{Name: "addresses",
			Type:     bigquery.RecordFieldType,
			Repeated: true,
			Schema: bigquery.Schema{
				{Name: "status", Type: bigquery.StringFieldType},
				{Name: "address", Type: bigquery.StringFieldType},
				{Name: "city", Type: bigquery.StringFieldType},
				{Name: "state", Type: bigquery.StringFieldType},
				{Name: "zip", Type: bigquery.StringFieldType},
				{Name: "numberOfYears", Type: bigquery.StringFieldType},
			}},
	}

	metaData := &bigquery.TableMetadata{
		Schema: sampleSchema,
	}
	tableRef := client.Dataset(datasetID).Table(tableID)
	if err := tableRef.Create(ctx, metaData); err != nil {
		return err
	}
	fmt.Fprintf(w, "created table %s\n", tableRef.FullyQualifiedName())
	return nil
}

// [END bigquery_nested_repeated_schema]
