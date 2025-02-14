// Copyright 2025 Google LLC
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

package viewdatasetaccesspolicy

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
	testfunctions "github.com/GoogleCloudPlatform/golang-samples/bigquery/bigquery_access_quickstart"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestViewDatasetAccessPolicies(t *testing.T) {
	tc := testutil.SystemTest(t)

	datasetName := "my_new_dataset_test"

	b := bytes.Buffer{}

	ctx := context.Background()

	// Creates bq client.
	client, err := testfunctions.TestClient(t, ctx)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}

	// Creates dataset.
	if err := client.Dataset(datasetName).Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
		t.Errorf("Failed to create dataset: %v", err)
	}

	// Once test is run, resources and clients are closed
	defer testfunctions.TestCleanup(t, ctx, client, datasetName)

	if err := viewDatasetAccessPolicies(&b, tc.ProjectID, datasetName); err != nil {
		t.Error(err)
	}

	if got, want := b.String(), "Details for Access entries in dataset"; !strings.Contains(got, want) {
		t.Errorf("viewDatasetAccessPolicies: expected %q to contain %q", got, want)
	}

}
