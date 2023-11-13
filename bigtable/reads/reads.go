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

// Package reads contains snippets related to reading data from Cloud Bigtable.
package reads

// [START bigtable_reads_imports]
import (
	"context"
	"fmt"
	"io"
	"strings"

	"cloud.google.com/go/bigtable"
)

// [END bigtable_reads_imports]

// [START bigtable_reads_row]
func readRow(w io.Writer, projectID, instanceID string, tableName string) error {
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
	rowKey := "phone#4c410523#20190501"
	row, err := tbl.ReadRow(ctx, rowKey)
	if err != nil {
		return fmt.Errorf("could not read row with key %s: %w", rowKey, err)
	}

	printRow(w, row)
	return nil
}

// [END bigtable_reads_row]

// [START bigtable_reads_row_partial]
func readRowPartial(w io.Writer, projectID, instanceID string, tableName string) error {
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
	rowKey := "phone#4c410523#20190501"
	row, err := tbl.ReadRow(ctx, rowKey, bigtable.RowFilter(bigtable.ColumnFilter("os_build")))
	if err != nil {
		return fmt.Errorf("could not read row with key %s: %w", rowKey, err)
	}

	printRow(w, row)

	return nil
}

// [END bigtable_reads_row_partial]

// [START bigtable_reads_rows]
func readRows(w io.Writer, projectID, instanceID string, tableName string) error {
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
	err = tbl.ReadRows(ctx, bigtable.RowList{"phone#4c410523#20190501", "phone#4c410523#20190502"},
		func(row bigtable.Row) bool {
			printRow(w, row)
			return true
		},
	)
	if err != nil {
		return fmt.Errorf("tbl.ReadRows: %w", err)
	}

	return nil
}

// [END bigtable_reads_rows]

// [START bigtable_reverse_scan]
func readRowsReversed(w io.Writer, projectID, instanceID string, tableName string) error {
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
	err = tbl.ReadRows(ctx, bigtable.NewRange("phone#5c10102", "phone#5c10103"),
		func(row bigtable.Row) bool {
			printRow(w, row)
			return true
		},
		bigtable.ReverseScan())
	if err != nil {
		return fmt.Errorf("tbl.ReadRows: %w", err)
	}

	return nil
}

// [END bigtable_reverse_scan]

// [START bigtable_reads_row_range]
func readRowRange(w io.Writer, projectID, instanceID string, tableName string) error {
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
	err = tbl.ReadRows(ctx, bigtable.NewRange("phone#4c410523#20190501", "phone#4c410523#201906201"),
		func(row bigtable.Row) bool {
			printRow(w, row)
			return true
		},
	)

	if err != nil {
		return fmt.Errorf("tbl.ReadRows: %w", err)
	}

	return nil
}

// [END bigtable_reads_row_range]

// [START bigtable_reads_row_ranges]
func readRowRanges(w io.Writer, projectID, instanceID string, tableName string) error {
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
	err = tbl.ReadRows(ctx, bigtable.RowRangeList{
		bigtable.NewRange("phone#4c410523#20190501", "phone#4c410523#201906201"),
		bigtable.NewRange("phone#5c10102#20190501", "phone#5c10102#201906201"),
	},
		func(row bigtable.Row) bool {
			printRow(w, row)
			return true
		},
	)
	if err != nil {
		return fmt.Errorf("tbl.ReadRows: %w", err)
	}

	return nil
}

// [END bigtable_reads_row_ranges]

// [START bigtable_reads_prefix]
func readPrefix(w io.Writer, projectID, instanceID string, tableName string) error {
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
	err = tbl.ReadRows(ctx, bigtable.PrefixRange("phone#"),
		func(row bigtable.Row) bool {
			printRow(w, row)
			return true
		},
	)
	if err != nil {
		return fmt.Errorf("tbl.ReadRows: %w", err)
	}

	return nil
}

// [END bigtable_reads_prefix]

// [START bigtable_reads_filter]
func readFilter(w io.Writer, projectID, instanceID string, tableName string) error {
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
	err = tbl.ReadRows(ctx, bigtable.RowRange{},
		func(row bigtable.Row) bool {
			printRow(w, row)
			return true
		},
		bigtable.RowFilter(bigtable.ValueFilter("PQ2A.*$")),
	)
	if err != nil {
		return fmt.Errorf("tbl.ReadRows: %w", err)
	}

	return nil
}

// [END bigtable_reads_filter]

// [START bigtable_reads_print]
func printRow(w io.Writer, row bigtable.Row) {
	fmt.Fprintf(w, "Reading data for %s:\n", row.Key())
	for columnFamily, cols := range row {
		fmt.Fprintf(w, "Column Family %s\n", columnFamily)
		for _, col := range cols {
			qualifier := col.Column[strings.IndexByte(col.Column, ':')+1:]
			fmt.Fprintf(w, "\t%s: %s @%d\n", qualifier, col.Value, col.Timestamp)
		}
	}
	fmt.Fprintln(w)
}

// [END bigtable_reads_print]
