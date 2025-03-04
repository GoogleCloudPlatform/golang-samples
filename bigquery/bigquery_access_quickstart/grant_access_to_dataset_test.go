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

package bigqueryaccessquickstart

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestGrantAccessDataset(t *testing.T) {

	tc := testutil.SystemTest(t)

	prefixer := testPrefix()
	prefix := fmt.Sprintf("%s_grant_access_to_dataset", prefixer)

	datasetName := fmt.Sprintf("%s_dataset", prefix)

	ctx := context.Background()

	var buff bytes.Buffer

	// Create BigQuery client.
	client, err := testClient(t)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	// Create dataset handler.
	dataset := client.Dataset(datasetName)
	defer testCleanup(t, client, datasetName)

	// Create dataset.
	if err := dataset.Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
		t.Fatalf("Failed to create dataset: %v", err)
	}

	if err := grantAccessToDataset(&buff, tc.ProjectID, datasetName); err != nil {
		t.Error(err)
	}

	if got, want := buff.String(), fmt.Sprintf("Details for Access entries in dataset %v.\n", datasetName); !strings.Contains(got, want) {
		t.Errorf("viewTableAccessPolicies: expected %q to contain %q", got, want)
	}
}
