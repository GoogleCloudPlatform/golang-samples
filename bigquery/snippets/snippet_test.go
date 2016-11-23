// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/net/context"
)

const ns = "snippets"

func TestAll(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}

	datasetID := testutil.NextDataset(ns)
	tableID := testutil.NextTable(ns)
	dstTableID := testutil.NextTable(ns)

	defer testutil.CleanupDatasets(t, datasetID)

	if err := createDataset(client, datasetID); err != nil {
		t.Errorf("failed to create dataset: %v", err)
	}
	if err := listDatasets(client); err != nil {
		t.Errorf("failed to create dataset: %v", err)
	}
	if err := createTable(client, datasetID, tableID); err != nil {
		t.Errorf("failed to create table: %v", err)
	}
	buf := &bytes.Buffer{}
	if err := listTables(client, buf, datasetID); err != nil {
		t.Errorf("failed to list tables: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, tableID) {
		t.Errorf("want table list %q to contain table %q", got, tableID)
	}
	if err := insertRows(client, datasetID, tableID); err != nil {
		t.Errorf("failed to insert rows: %v", err)
	}
	if err := listRows(client, datasetID, tableID); err != nil {
		t.Errorf("failed to list rows: %v", err)
	}
	if err := browseTable(client, datasetID, tableID); err != nil {
		t.Errorf("failed to list rows: %v", err)
	}
	if err := asyncQuery(client, datasetID, tableID); err != nil {
		t.Errorf("failed to async query: %v", err)
	}
	if err := copyTable(client, datasetID, tableID, dstTableID); err != nil {
		t.Errorf("failed to copy table: %v", err)
	}
	if err := deleteTable(client, datasetID, tableID); err != nil {
		t.Errorf("failed to delete table: %v", err)
	}
	if err := deleteTable(client, datasetID, dstTableID); err != nil {
		t.Errorf("failed to delete table: %v", err)
	}
}

func TestImportExport(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}

	datasetID := testutil.NextDataset("importexport")
	tableID := testutil.NextTable("importexport")

	if err := createDataset(client, datasetID); err != nil {
		t.Errorf("failed to create dataset: %v", err)
	}
	defer testutil.CleanupDatasets(t, datasetID)

	schema := bigquery.Schema{
		&bigquery.FieldSchema{Name: "Year", Type: bigquery.IntegerFieldType},
		&bigquery.FieldSchema{Name: "City", Type: bigquery.StringFieldType},
	}
	if err := client.Dataset(datasetID).Table(tableID).Create(ctx, schema); err != nil {
		t.Errorf("failed to create dataset: %v", err)
	}
	if err := importFromFile(client, datasetID, tableID, "testdata/olympics.csv"); err != nil {
		t.Fatalf("failed to import from file: %v", err)
	}

	bucket := testutil.NextBucket("importexport")
	const object = "values.csv"

	testutil.BucketsMustExist(t, bucket)
	defer testutil.CleanupBuckets(t, bucket)

	gcsURI := fmt.Sprintf("gs://%s/%s", bucket, object)
	if err := exportToGCS(client, datasetID, tableID, gcsURI); err != nil {
		t.Errorf("failed to export to %v: %v", gcsURI, err)
	}
}
