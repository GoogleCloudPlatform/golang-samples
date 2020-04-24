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

// Package table demonstrates interactions with BigQuery Table resources,
// which included managed tables, federation, and logical views.
package table

import (
	"context"
	"io/ioutil"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/bigquery/snippets/bqtestutil"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestTables(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	meta := &bigquery.DatasetMetadata{
		Location: "US", // See https://cloud.google.com/bigquery/docs/locations
	}
	testDatasetID, err := bqtestutil.UniqueBQName("snippet_table_tests")
	if err != nil {
		t.Fatalf("couldn't generate unique resource name: %v", err)
	}
	if err := client.Dataset(testDatasetID).Create(ctx, meta); err != nil {
		t.Fatalf("failed to create test dataset: %v", err)
	}
	// Cleanup dataset at end of test.
	defer client.Dataset(testDatasetID).DeleteWithContents(ctx)

	testDatasetID2, err := bqtestutil.UniqueBQName("second_snippet_table_tests")
	if err != nil {
		t.Fatalf("couldn't generate unique resource name: %v", err)
	}
	if err := client.Dataset(testDatasetID2).Create(ctx, meta); err != nil {
		t.Fatalf("failed to create test dataset: %v", err)
	}
	// Cleanup dataset at end of test.
	defer client.Dataset(testDatasetID2).DeleteWithContents(ctx)

	testTableID, err := bqtestutil.UniqueBQName("testtable")
	if err != nil {
		t.Fatalf("couldn't generate unique table id: %v", err)
	}
	if err := createTableExplicitSchema(tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("createTableExplicitSchema(%q %q): %v", testDatasetID, testTableID, err)
	}
	if err := insertRows(tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("insertRows(%q %q): %v", testDatasetID, testTableID, err)
	}
	if err := printTableInfo(ioutil.Discard, tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("printTableInfo(%q %q): %v", testDatasetID, testTableID, err)
	}
	if err := browseTable(ioutil.Discard, tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("browseTable(%q %q): %v", testDatasetID, testTableID, err)
	}
	if err := deleteAndUndeleteTable(tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("deleteAndUndeleteTable(%q %q): %v", testDatasetID, testTableID, err)
	}

	testTableID, err = bqtestutil.UniqueBQName("testtable")
	if err != nil {
		t.Fatalf("couldn't generate unique table id: %v", err)
	}
	if err := createTableComplexSchema(ioutil.Discard, tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("createTableComplexSchema(%q %q): %v", testDatasetID, testTableID, err)
	}
	if err := updateTableDescription(tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("updateTableDescription(%q %q): %v", testDatasetID, testTableID, err)
	}
	if err := updateTableExpiration(tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("updateTableExpiration(%q %q): %v", testDatasetID, testTableID, err)
	}

	testTableID, err = bqtestutil.UniqueBQName("testtable")
	if err != nil {
		t.Fatalf("couldn't generate unique table id: %v", err)
	}
	if err := createTablePartitioned(tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("createTablePartitioned(%q %q): %v", testDatasetID, testTableID, err)
	}

	testTableID, err = bqtestutil.UniqueBQName("testtable")
	if err != nil {
		t.Fatalf("couldn't generate unique table id: %v", err)
	}
	if err := createTableClustered(tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("createTableClustered(%q %q): %v", testDatasetID, testTableID, err)
	}

	testTableID, err = bqtestutil.UniqueBQName("testtable")
	if err != nil {
		t.Fatalf("couldn't generate unique table id: %v", err)
	}

	t.Run("cmektests", func(t *testing.T) {
		if bqtestutil.RunCMEKTests() {
			t.Skip("skipping CMEK tests")
		}
		if err := createTableWithCMEK(tc.ProjectID, testDatasetID, testTableID); err != nil {
			t.Fatalf("createTableWithCMEK(%q %q): %v", testDatasetID, testTableID, err)
		}
		if err := updateTableChangeCMEK(tc.ProjectID, testDatasetID, testTableID); err != nil {
			t.Fatalf("updateTableChangeCMEK(%q %q): %v", testDatasetID, testTableID, err)
		}
	})

	testTableID, err = bqtestutil.UniqueBQName("testtable")
	if err != nil {
		t.Fatalf("couldn't generate unique table id: %v", err)
	}
	if err := createView(tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("createView(%q %q): %v", testDatasetID, testTableID, err)
	}
	if err := getView(ioutil.Discard, tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("getView(%q %q): %v", testDatasetID, testTableID, err)
	}
	if err := updateViewDelegated(tc.ProjectID, testDatasetID2, testDatasetID, testTableID); err != nil {
		t.Fatalf("updateViewDelegated(%q %q): %v", testDatasetID, testTableID, err)
	}

	testTableID, err = bqtestutil.UniqueBQName("testtable")
	if err != nil {
		t.Fatalf("couldn't generate unique table id: %v", err)
	}
	if err := relaxTableAPI(tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("relaxTableAPI(%q %q): %v", testDatasetID, testTableID, err)
	}
	if err := updateTableAddColumn(tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("updateTableAddColumn(%q %q): %v", testDatasetID, testTableID, err)
	}
	if err := addTableLabel(tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("addTableLabel(%q %q): %v", testDatasetID, testTableID, err)
	}
	if err := tableLabels(ioutil.Discard, tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("tableLabels(%q %q): %v", testDatasetID, testTableID, err)
	}
	if err := deleteTableLabel(tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("deleteTableLabel(%q %q): %v", testDatasetID, testTableID, err)
	}
	if err := deleteTable(tc.ProjectID, testDatasetID, testTableID); err != nil {
		t.Fatalf("deleteTable(%q %q): %v", testDatasetID, testTableID, err)
	}

	if err := listTables(ioutil.Discard, tc.ProjectID, testDatasetID); err != nil {
		t.Fatalf("deleteTable(%q): %v", testDatasetID, err)
	}

}
