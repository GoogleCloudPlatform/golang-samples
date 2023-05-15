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

// [START bigquery_relax_column]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// relaxTableAPI demonstrates modifying the schema of a table to remove the requirement that columns allow
// no NULL values.
func relaxTableAPI(projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydatasetid"
	// tableID := "mytableid"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	// Setup: We first create a table with a schema that's restricts NULL values.
	sampleSchema := bigquery.Schema{
		{Name: "full_name", Type: bigquery.StringFieldType, Required: true},
		{Name: "age", Type: bigquery.IntegerFieldType, Required: true},
	}
	original := &bigquery.TableMetadata{
		Schema: sampleSchema,
	}
	if err := client.Dataset(datasetID).Table(tableID).Create(ctx, original); err != nil {
		return err
	}

	tableRef := client.Dataset(datasetID).Table(tableID)
	meta, err := tableRef.Metadata(ctx)
	if err != nil {
		return err
	}
	// Iterate through the schema to set all Required fields to false (nullable).
	var relaxed bigquery.Schema
	for _, v := range meta.Schema {
		v.Required = false
		relaxed = append(relaxed, v)
	}
	newMeta := bigquery.TableMetadataToUpdate{
		Schema: relaxed,
	}
	if _, err := tableRef.Update(ctx, newMeta, meta.ETag); err != nil {
		return err
	}
	return nil
}

// [END bigquery_relax_column]
