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

// Package dataset demonstrates interactions with BigQuery's dataset resources.
// Examples include lifecycle operations such as creation, modification, and
// deletion.
package dataset

import (
	"bytes"
	"context"
	"io/ioutil"
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/bigquery/snippets/bqtestutil"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDatasets(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	datasetID, err := bqtestutil.UniqueBQName("golang_snippettest_dataset")
	if err != nil {
		t.Fatal(err)
	}

	if err := createDataset(tc.ProjectID, datasetID); err != nil {
		t.Errorf("createDataset(%q): %v", datasetID, err)
	}
	// Cleanup dataset at end of test.
	defer client.Dataset(datasetID).DeleteWithContents(ctx)

	if err := updateDatasetAccessControl(tc.ProjectID, datasetID); err != nil {
		t.Errorf("updateDataSetAccessControl(%q): %v", datasetID, err)
	}
	if err := revokeDatasetAccess(tc.ProjectID, datasetID, "sample.bigquery.dev@gmail.com"); err != nil {
		t.Errorf("revokeDatasetAccess(%q): %v", datasetID, err)
	}
	if err := addDatasetLabel(tc.ProjectID, datasetID); err != nil {
		t.Errorf("updateDatasetAddLabel: %v", err)
	}

	buf := &bytes.Buffer{}
	if err := printDatasetLabels(buf, tc.ProjectID, datasetID); err != nil {
		t.Errorf("printDatasetLabels(%q): %v", datasetID, err)
	}
	want := "color:green"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("getDatasetLabel(%q) expected %q to contain %q", datasetID, got, want)
	}

	if err := addDatasetLabel(tc.ProjectID, datasetID); err != nil {
		t.Errorf("updateDatasetAddLabel: %v", err)
	}
	buf.Reset()
	if err := listDatasetsByLabel(buf, tc.ProjectID); err != nil {
		t.Errorf("listDatasetsByLabel: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, datasetID) {
		t.Errorf("listDatasetsByLabel expected %q to contain %q", got, want)
	}
	if err := deleteDatasetLabel(tc.ProjectID, datasetID); err != nil {
		t.Errorf("updateDatasetDeleteLabel: %v", err)
	}
	if err := printDatasetInfo(ioutil.Discard, tc.ProjectID, datasetID); err != nil {
		t.Errorf("printDatasetInfo: %v", err)
	}

	// Test empty dataset creation/ttl/delete.
	deletionDatasetID, err := bqtestutil.UniqueBQName("golang_example_quickdelete")
	if err != nil {
		t.Fatal(err)
	}
	if err := createDataset(tc.ProjectID, deletionDatasetID); err != nil {
		t.Errorf("createDataset(%q): %v", deletionDatasetID, err)
	}
	if err = updateDatasetDefaultExpiration(tc.ProjectID, deletionDatasetID); err != nil {
		t.Errorf("updateDatasetDefaultExpiration(%q): %v", deletionDatasetID, err)
	}
	if err := deleteDataset(tc.ProjectID, deletionDatasetID); err != nil {
		t.Errorf("deleteEmptyDataset(%q): %v", deletionDatasetID, err)
	}

	if err := updateDatasetDescription(tc.ProjectID, datasetID); err != nil {
		t.Errorf("updateDatasetDescription(%q): %v", datasetID, err)
	}
	if err := listDatasets(tc.ProjectID, ioutil.Discard); err != nil {
		t.Errorf("listDatasets: %v", err)
	}

}
