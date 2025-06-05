// Copyright 2025 Google LLC
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

package deletes

// [START bigtable_streaming_and_batching_asyncio]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigtable"
)

// streamingAndBatching starts a stream of data (reading rows), batches them, and then goes through the batch and deletes all the cells in column
func streamingAndBatching(w io.Writer, projectID, instanceID string, tableName string) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance-id"
	// tableName := "mobile-time-series"

	ctx := context.Background()
	client, err := bigtable.NewClient(ctx, projectID, instanceID)
	if err != nil {
		return fmt.Errorf("bigtable.NewClient: %w", err)
	}
	defer client.Close()
	tbl := client.Open(tableName)

	// Slices to hold the row keys and the corresponding mutations.
	var rowKeys []string
	var mutations []*bigtable.Mutation

	// Read all rows from the table.
	err = tbl.ReadRows(ctx, bigtable.InfiniteRange(""), func(row bigtable.Row) bool {
		// For each row, create a mutation to delete the specified cell.
		mut := bigtable.NewMutation()
		mut.DeleteCellsInColumn("cell_plan", "data_plan_01gb")

		// Append the row key and mutation to the slices.
		rowKeys = append(rowKeys, row.Key())
		mutations = append(mutations, mut)

		// Continue processing rows.
		return true
	})
	if err != nil {
		return fmt.Errorf("tbl.ReadRows: %w", err)
	}

	if len(mutations) == 0 {
		return nil
	}
	// If there are mutations to apply, apply them in a single bulk request.
	// ApplyBulk returns a slice of errors, one for each mutation.
	var errs []error
	if errs, err = tbl.ApplyBulk(ctx, rowKeys, mutations); err != nil {
		return fmt.Errorf("tbl.ApplyBulk: %w", err)
	}
	if errs != nil {
		// Log any individual errors that occurred during the bulk operation.
		var errorCount int
		for _, individualErr := range errs {
			if individualErr != nil {
				fmt.Fprintf(w, "Error applying mutation: %v\n", individualErr)
				errorCount++
			}
		}
		if errorCount > 0 {
			return fmt.Errorf("encountered %d error(s) out of %d mutations", errorCount, len(errs))
		}
	}

	fmt.Fprintf(w, "Successfully deleted cells from all rows")
	return nil
}

// [END bigtable_streaming_and_batching_asyncio]
