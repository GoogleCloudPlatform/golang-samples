// Copyright 2024 Google LLC
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

// Package writes contains snippets related to writing to Cloud Bigtable.
package writes

// [START bigtable_writes_aggregate]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/bigtable"
)

func writeAggregate(w io.Writer, projectID, instanceID string, tableName string) error {
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
	columnFamilyName := "view_count"
	viewTimestamp, err := time.Parse(time.RFC3339, "2024-03-13T12:41:34Z")
	if err != nil {
		return err
	}
	hourlyBucket := viewTimestamp.Truncate(time.Hour)

	mut := bigtable.NewMutation()
	mut.AddIntToCell(columnFamilyName, "views", bigtable.Time(hourlyBucket), 1)

	rowKey := "page#index.html"
	if err := tbl.Apply(ctx, rowKey, mut); err != nil {
		return fmt.Errorf("Apply: %w", err)
	}

	fmt.Fprintf(w, "Successfully wrote row: %s\n", rowKey)
	return nil
}

// [END bigtable_writes_aggregate]
