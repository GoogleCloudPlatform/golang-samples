// Copyright 2021 Google LLC
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

// [START bigquery_create_table_from_template]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// createTableFromTemplateTable demonstrates how to use the properties of one
// table (schema, partitioning, clustering) to create a new empty table with
// the same configuration.
func createTableFromTemplateTable(srcProjectID, srcDatasetID, srcTableID, dstProjectID, dstDatasetID, dstTableID string) error {
	// srcProjectID := "bigquery-public-data"
	// srcDatasetID := "samples"
	// srcTableID := "shakespeare"
	// dstProjectID := "my-project-id"
	// dstDatasetID := "mydataset"
	// dstTableID := "mytable"
	ctx := context.Background()

	// We'll construct the client based on the destination project.
	client, err := bigquery.NewClient(ctx, dstProjectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	srcTableRef := client.DatasetInProject(srcProjectID, srcDatasetID).Table(srcTableID)

	srcMeta, err := srcTableRef.Metadata(ctx)
	if err != nil {
		return fmt.Errorf("failed to get source table metadata: %w", err)
	}

	dstTableRef := client.Dataset(dstDatasetID).Table(dstTableID)

	// We'll use some (but not all) of the metadata from the source table
	// to define the destination table.  Other properties to consider include
	// attributes like expiration policy and managed encryption settings.
	dstMeta := &bigquery.TableMetadata{
		Description: fmt.Sprintf("table structure copied from %s.%s.%s",
			srcTableRef.ProjectID, srcTableRef.DatasetID, srcTableRef.TableID),
		Schema:     srcMeta.Schema,
		Clustering: srcMeta.Clustering,
		// A table will only have one partitioning configuration, the others will be empty/nil.
		TimePartitioning:  srcMeta.TimePartitioning,
		RangePartitioning: srcMeta.RangePartitioning,
		// Retain the same enforcement of partition filtering as well.
		RequirePartitionFilter: srcMeta.RequirePartitionFilter,
	}

	if err := dstTableRef.Create(ctx, dstMeta); err != nil {
		return err
	}

	return nil
}

// [END bigquery_create_table_from_template]
