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

package dataset

import (
	"bytes"
	"context"
	"io/ioutil"
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
	bqtestutil "github.com/GoogleCloudPlatform/golang-samples/bigquery/snippets/testutil"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestAll(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}

	datasetID := bqtestutil.UniqueBQName("golang_snippettest_dataset")
	if err := createDataset(client, datasetID); err != nil {
		t.Errorf("createDataset(%q): %v", datasetID, err)
	}
	// Cleanup dataset at end of test.
	defer client.Dataset(datasetID).DeleteWithContents(ctx)

	if err := updateDatasetAccessControl(client, datasetID); err != nil {
		t.Errorf("updateDataSetAccessControl(%q): %v", datasetID, err)
	}
	if err := addDatasetLabel(client, datasetID); err != nil {
		t.Errorf("updateDatasetAddLabel: %v", err)
	}

	buf := &bytes.Buffer{}
	if err := printDatasetLabels(client, buf, datasetID); err != nil {
		t.Errorf("printDatasetLabels(%q): %v", datasetID, err)
	}
	want := "color:green"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("getDatasetLabel(%q) expected %q to contain %q", datasetID, got, want)
	}

	if err := addDatasetLabel(client, datasetID); err != nil {
		t.Errorf("updateDatasetAddLabel: %v", err)
	}
	buf.Reset()
	if err := listDatasetsByLabel(client, buf); err != nil {
		t.Errorf("listDatasetsByLabel: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, datasetID) {
		t.Errorf("listDatasetsByLabel expected %q to contain %q", got, want)
	}
	if err := deleteDatasetLabel(client, datasetID); err != nil {
		t.Errorf("updateDatasetDeleteLabel: %v", err)
	}
	if err := printDatasetInfo(client, ioutil.Discard, datasetID); err != nil {
		t.Errorf("printDatasetInfo: %v", err)
	}

	// Test empty dataset creation/ttl/delete.
	deletionDatasetID := bqtestutil.UniqueBQName("golang_example_quickdelete")
	if err := createDataset(client, deletionDatasetID); err != nil {
		t.Errorf("createDataset(%q): %v", deletionDatasetID, err)
	}
	if err = updateDatasetDefaultExpiration(client, deletionDatasetID); err != nil {
		t.Errorf("updateDatasetDefaultExpiration(%q): %v", deletionDatasetID, err)
	}
	if err := deleteDataset(client, deletionDatasetID); err != nil {
		t.Errorf("deleteEmptyDataset(%q): %v", deletionDatasetID, err)
	}

	if err := updateDatasetDescription(client, datasetID); err != nil {
		t.Errorf("updateDatasetDescription(%q): %v", datasetID, err)
	}
	if err := listDatasets(client, ioutil.Discard); err != nil {
		t.Errorf("listDatasets: %v", err)
	}

}
