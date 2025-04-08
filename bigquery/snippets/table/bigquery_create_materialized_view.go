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

// [START bigquery_create_materialized_view]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// createMaterializedView demonstrates generated a materialized view based on an existing
// base table.
func createMaterializedView(projectID, datasetID, baseTableID, viewID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydatasetid"
	// baseTableID := "mytableid"
	// viewID := "myviewid"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// Get an appropriately escaped table identifier suitable for use in a standard SQL query.
	tableStr, err := client.Dataset(datasetID).Table(baseTableID).Identifier(bigquery.StandardSQLID)
	if err != nil {
		return fmt.Errorf("couldn't construct identifier: %w", err)
	}

	metaData := &bigquery.TableMetadata{
		MaterializedView: &bigquery.MaterializedViewDefinition{
			Query: fmt.Sprintf(`SELECT MAX(TimestampField) AS TimestampField, StringField, 
					  MAX(BooleanField) AS BooleanField FROM %s GROUP BY StringField`, tableStr),
		}}

	viewRef := client.Dataset(datasetID).Table(viewID)
	if err := viewRef.Create(ctx, metaData); err != nil {
		return err
	}
	return nil
}

// [END bigquery_create_materialized_view]
