// Copyright 2026 Google LLC
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

// [START bigtable_delete_from_row]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigtable"
)

func deleteFromRow(w io.Writer, projectID, instanceID, tableName string) error {
	ctx := context.Background()
	client, err := bigtable.NewClient(ctx, projectID, instanceID)
	if err != nil {
		return fmt.Errorf("bigtable.NewClient: %w", err)
	}
	defer client.Close()
	tbl := client.Open(tableName)

	// Use a specific row key that exists in your table.
	rowKey := "phone#4c410523#20190501"
	mut := bigtable.NewMutation()
	mut.DeleteRow()

	if err := tbl.Apply(ctx, rowKey, mut); err != nil {
		return fmt.Errorf("tbl.Apply: %w", err)
	}

	fmt.Fprintf(w, "Successfully deleted row: %s\n", rowKey)
	return nil
}

// [END bigtable_delete_from_row]
