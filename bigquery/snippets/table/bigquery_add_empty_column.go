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

// [START bigquery_add_empty_column]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// updateTableAddColumn demonstrates modifying the schema of a table to append an additional column.
func updateTableAddColumn(projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	tableRef := client.Dataset(datasetID).Table(tableID)
	meta, err := tableRef.Metadata(ctx)
	if err != nil {
		return err
	}
	newSchema := append(meta.Schema,
		&bigquery.FieldSchema{Name: "phone", Type: bigquery.StringFieldType},
	)
	update := bigquery.TableMetadataToUpdate{
		Schema: newSchema,
	}
	if _, err := tableRef.Update(ctx, update, meta.ETag); err != nil {
		return err
	}
	return nil
}

// [END bigquery_add_empty_column]
