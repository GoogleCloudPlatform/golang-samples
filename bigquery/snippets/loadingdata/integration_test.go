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

// Package loadingdata demonstrates interactions with BigQuery's batch
// ingestion mechanism.

package loadingdata

import (
	"context"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/bigquery/snippets/bqtestutil"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestImportSnippets(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	// Control a job lifecycle explicitly: create, report status, cancel.
	testDatasetID, err := bqtestutil.UniqueBQName("golang_snippets_loading")
	if err != nil {
		t.Fatalf("couldn't generate unique resource name: %v", err)
	}

	meta := &bigquery.DatasetMetadata{
		Location: "US", // See https://cloud.google.com/bigquery/docs/locations
	}
	if err := client.Dataset(testDatasetID).Create(ctx, meta); err != nil {
		t.Fatalf("failed to create test dataset: %v", err)
	}
	// Cleanup dataset at end of test.
	defer client.Dataset(testDatasetID).DeleteWithContents(ctx)

	// Run these in a group.  They're batch workloads and can run concurrently.
	t.Run("group", func(t *testing.T) {
		t.Run("importCSVFromFile", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_load_from_file"
			filename := "../testdata/people.csv"
			if err := importCSVFromFile(tc.ProjectID, testDatasetID, tableID, filename); err != nil {
				t.Errorf("importCSVFromFile(%q): %v", testDatasetID, err)
			}
		})
		t.Run("importCSVAutodetectSchema", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_load_table_gcs_csv_autodetect"
			if err := importCSVAutodetectSchema(tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("importCSVAutodetectSchema(%q): %v", testDatasetID, err)
			}
		})
		t.Run("importCSVTruncate", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_load_table_gcs_csv_truncate"
			if err := importCSVTruncate(tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("importCSVTruncate(%q): %v", testDatasetID, err)
			}
		})
		t.Run("importCSVExplicitSchema", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_load_table_gcs_csv"
			if err := importCSVExplicitSchema(tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("importCSVExplicitSchema(%q): %v", testDatasetID, err)
			}
		})
		t.Run("importJSONExplicitSchema", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_load_table_gcs_json"
			if err := importJSONExplicitSchema(tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("importJSONExplicitSchema(%q): %v", testDatasetID, err)
			}
		})
		t.Run("importJSONAutodetectSchema", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_load_table_gcs_json_autodetect"
			if err := importJSONAutodetectSchema(tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("importJSONAutodetectSchema(%q): %v", testDatasetID, err)
			}
		})
		t.Run("importJSONWithCMEK", func(t *testing.T) {
			if bqtestutil.RunCMEKTests() {
				t.Skip("skipping CMEK tests")
			}
			t.Parallel()
			tableID := "bigquery_load_table_gcs_json_cmek"
			if err := importJSONWithCMEK(tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("importJSONWithCMEK(%q): %v", testDatasetID, err)
			}
		})
		t.Run("importJSONTruncate", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_load_table_gcs_json_truncate"
			if err := importJSONTruncate(tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("importJSONTruncate(%q): %v", testDatasetID, err)
			}
		})
		t.Run("importORC", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_load_table_gcs_orc"
			if err := importORC(tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("importORC(%q): %v", testDatasetID, err)
			}
		})
		t.Run("importORCTruncate", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_load_table_gcs_orc_truncate"
			if err := importORCTruncate(tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("importORCTruncate(%q): %v", testDatasetID, err)
			}
		})
		t.Run("importParquet", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_load_table_gcs_parquet"
			if err := importParquet(tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("importParquet(%q): %v", testDatasetID, err)
			}
		})
		t.Run("importParquetTruncate", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_load_table_gcs_parquet_truncate"
			if err := importParquetTruncate(tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("importParquetTruncate(%q): %v", testDatasetID, err)
			}
		})
		t.Run("createTableAndWidenLoad", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_add_column_load_append"
			filename := "../testdata/people.csv"
			if err := createTableAndWidenLoad(tc.ProjectID, testDatasetID, tableID, filename); err != nil {
				t.Errorf("createTableAndWidenLoad(%q): %v", testDatasetID, err)
			}
		})
		t.Run("relaxTableImport", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_relax_column_load_append"
			filename := "../testdata/people.csv"
			if err := relaxTableImport(tc.ProjectID, testDatasetID, tableID, filename); err != nil {
				t.Errorf("relaxTableImport(%q): %v", testDatasetID, err)
			}
		})
		t.Run("importPartitionedTable", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_load_table_partitioned"
			if err := importPartitionedTable(tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("importPartitionedTable(%q): %v", testDatasetID, err)
			}
		})
		t.Run("importClusteredSampleTable", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_load_table_clustered"
			if err := importClusteredTable(tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("importClusteredTable(%q): %v", testDatasetID, err)
			}
		})
	})

}
