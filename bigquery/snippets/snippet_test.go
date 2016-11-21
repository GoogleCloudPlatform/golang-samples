// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/net/context"
)

func TestAll(t *testing.T) {
	tc := testutil.SystemTest(t)

	client, err := bigquery.NewClient(context.Background(), tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	datasetID := fmt.Sprintf("golang_example_dataset_%d", time.Now().Unix())
	if err := createDataset(client, datasetID); err != nil {
		t.Errorf("failed to create dataset: %v", err)
	}
	if err := listDatasets(client); err != nil {
		t.Errorf("failed to create dataset: %v", err)
	}

	tableID := fmt.Sprintf("golang_example_table_%d", time.Now().Unix())
	if err := createTable(client, datasetID, tableID); err != nil {
		t.Errorf("failed to create table: %v", err)
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

	dstTableID := fmt.Sprintf("golang_example_tabledst_%d", time.Now().Unix())
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
