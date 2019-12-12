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

// Package filters contains snippets related to reading data from Cloud Bigtable
// with various filters.
package filters

// [START bigtable_filters_limit_row_sample]
// [START bigtable_filters_limit_row_regex]
// [START bigtable_filters_limit_cells_per_col]
// [START bigtable_filters_limit_cells_per_row]
// [START bigtable_filters_limit_cells_per_row_offset]
// [START bigtable_filters_limit_col_family_regex]
// [START bigtable_filters_limit_col_qualifier_regex]
// [START bigtable_filters_limit_col_range]
// [START bigtable_filters_limit_value_range]
// [START bigtable_filters_limit_value_regex]
// [START bigtable_filters_limit_timestamp_range]
// [START bigtable_filters_limit_block_all]
// [START bigtable_filters_limit_pass_all]
// [START bigtable_filters_modify_strip_value]
// [START bigtable_filters_modify_apply_label]
// [START bigtable_filters_composing_chain]
// [START bigtable_filters_composing_interleave]
// [START bigtable_filters_composing_condition]
import (
	"context"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"

	"cloud.google.com/go/bigtable"
)

// [END bigtable_filters_limit_row_sample]
// [END bigtable_filters_limit_row_regex]
// [END bigtable_filters_limit_cells_per_col]
// [END bigtable_filters_limit_cells_per_row]
// [END bigtable_filters_limit_cells_per_row_offset]
// [END bigtable_filters_limit_col_family_regex]
// [END bigtable_filters_limit_col_qualifier_regex]
// [END bigtable_filters_limit_col_range]
// [END bigtable_filters_limit_value_range]
// [END bigtable_filters_limit_value_regex]
// [END bigtable_filters_limit_timestamp_range]
// [END bigtable_filters_limit_block_all]
// [END bigtable_filters_limit_pass_all]
// [END bigtable_filters_modify_strip_value]
// [END bigtable_filters_modify_apply_label]
// [END bigtable_filters_composing_chain]
// [END bigtable_filters_composing_interleave]
// [END bigtable_filters_composing_condition]

// [START bigtable_filters_limit_row_sample]
func filterLimitRowSample(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.RowSampleFilter(.75)
	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_limit_row_sample]
// [START bigtable_filters_limit_row_regex]
func filterLimitRowRegex(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.RowKeyFilter(".*#20190501$")
	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_limit_row_regex]
// [START bigtable_filters_limit_cells_per_col]
func filterLimitCellsPerCol(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.LatestNFilter(2)
	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_limit_cells_per_col]
// [START bigtable_filters_limit_cells_per_row]
func filterLimitCellsPerRow(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.CellsPerRowLimitFilter(2)
	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_limit_cells_per_row]
// [START bigtable_filters_limit_cells_per_row_offset]
func filterLimitCellsPerRowOffset(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.CellsPerRowOffsetFilter(2)
	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_limit_cells_per_row_offset]
// [START bigtable_filters_limit_col_family_regex]
func filterLimitColFamilyRegex(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.FamilyFilter("stats_.*$")
	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_limit_col_family_regex]
// [START bigtable_filters_limit_col_qualifier_regex]
func filterLimitColQualifierRegex(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.ColumnFilter("connected_.*$")
	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_limit_col_qualifier_regex]
// [START bigtable_filters_limit_col_range]
func filterLimitColRange(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.ColumnRangeFilter("cell_plan", "data_plan_01gb", "data_plan_10gb")
	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_limit_col_range]
// [START bigtable_filters_limit_value_range]
func filterLimitValueRange(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.ValueRangeFilter([]byte("PQ2A.190405"), []byte("PQ2A.190406"))
	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_limit_value_range]
// [START bigtable_filters_limit_value_regex]
func filterLimitValueRegex(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.ValueFilter("PQ2A.*$")
	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_limit_value_regex]
// [START bigtable_filters_limit_timestamp_range]
func filterLimitTimestampRange(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.TimestampRangeFilter(
		(bigtable.Timestamp(0)).Time(),
		bigtable.Timestamp(bigtable.Now().TruncateToMilliseconds() - 60*60*1000*1000).Time())

	//filter := bigtable.TimestampRangeFilterMicros()
	//filter := bigtable.ValueFilter("PQ2A.*$")

	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_limit_timestamp_range]
// [START bigtable_filters_limit_block_all]
func filterLimitBlockAll(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.BlockAllFilter()
	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_limit_block_all]
// [START bigtable_filters_limit_pass_all]
func filterLimitPassAll(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.PassAllFilter()
	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_limit_pass_all]
// [START bigtable_filters_modify_strip_value]
func filterModifyStripValue(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.StripValueFilter()
	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_modify_strip_value]

// [START bigtable_filters_composing_chain]
func filterComposingChain(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.ChainFilters(bigtable.LatestNFilter(1), bigtable.FamilyFilter("cell_plan"))
	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_composing_chain]
// [START bigtable_filters_composing_interleave]
func filterComposingInterleave(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.InterleaveFilters(
		bigtable.ValueFilter("true"),
		bigtable.ColumnFilter("os_build"))
	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_composing_interleave]
// [START bigtable_filters_composing_condition]
func filterComposingCondition(w io.Writer, projectID, instanceID string, tableName string) error {
	filter := bigtable.ConditionFilter(
		bigtable.ChainFilters(
			bigtable.ColumnFilter("data_plan_10gb"),
			bigtable.ValueFilter("true")),
		bigtable.StripValueFilter(),
		bigtable.PassAllFilter())
	readFilter(w, projectID, instanceID, tableName, filter)
	return nil
}

// [END bigtable_filters_composing_condition]

// [START bigtable_filters_limit_row_sample]
// [START bigtable_filters_limit_row_regex]
// [START bigtable_filters_limit_cells_per_col]
// [START bigtable_filters_limit_cells_per_row]
// [START bigtable_filters_limit_cells_per_row_offset]
// [START bigtable_filters_limit_col_family_regex]
// [START bigtable_filters_limit_col_qualifier_regex]
// [START bigtable_filters_limit_col_range]
// [START bigtable_filters_limit_value_range]
// [START bigtable_filters_limit_value_regex]
// [START bigtable_filters_limit_timestamp_range]
// [START bigtable_filters_limit_block_all]
// [START bigtable_filters_limit_pass_all]
// [START bigtable_filters_modify_strip_value]
// [START bigtable_filters_modify_apply_label]
// [START bigtable_filters_composing_chain]
// [START bigtable_filters_composing_interleave]
// [START bigtable_filters_composing_condition]
func readFilter(w io.Writer, projectID, instanceID string, tableName string, filter bigtable.Filter) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance-id"
	// tableName := "mobile-time-series"

	ctx := context.Background()
	client, err := bigtable.NewClient(ctx, projectID, instanceID)
	if err != nil {
		return fmt.Errorf("bigtable.NewAdminClient: %v", err)
	}
	tbl := client.Open(tableName)
	err = tbl.ReadRows(ctx, bigtable.RowRange{},
		func(row bigtable.Row) bool {
			printRow(w, row)
			fmt.Fprintf(w, "\n")
			return true
		}, bigtable.RowFilter(filter))

	if err = client.Close(); err != nil {
		log.Fatalf("Could not close data operations client: %v", err)
	}

	return nil
}

func printRow(w io.Writer, row bigtable.Row) {
	fmt.Fprintf(w, "Reading data for %s:\n", row.Key())
	var keys []string
	for k := range row {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, columnFamily := range keys {
		cols := row[columnFamily]
		fmt.Fprintf(w, "Column Family %s\n", columnFamily)

		for _, col := range cols {
			qualifier := col.Column[strings.IndexByte(col.Column, ':')+1:]
			//labels := col.
			fmt.Fprintf(w, "\t%s: %s @%d\n", qualifier, col.Value, col.Timestamp)
		}
	}
}

// [END bigtable_filters_limit_row_sample]
// [END bigtable_filters_limit_row_regex]
// [END bigtable_filters_limit_cells_per_col]
// [END bigtable_filters_limit_cells_per_row]
// [END bigtable_filters_limit_cells_per_row_offset]
// [END bigtable_filters_limit_col_family_regex]
// [END bigtable_filters_limit_col_qualifier_regex]
// [END bigtable_filters_limit_col_range]
// [END bigtable_filters_limit_value_range]
// [END bigtable_filters_limit_value_regex]
// [END bigtable_filters_limit_timestamp_range]
// [END bigtable_filters_limit_block_all]
// [END bigtable_filters_limit_pass_all]
// [END bigtable_filters_modify_strip_value]
// [END bigtable_filters_modify_apply_label]
// [END bigtable_filters_composing_chain]
// [END bigtable_filters_composing_interleave]
// [END bigtable_filters_composing_condition]
