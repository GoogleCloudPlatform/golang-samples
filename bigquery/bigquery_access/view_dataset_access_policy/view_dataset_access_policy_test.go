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
	"log"
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestViewDatasetAccessPolicies(t *testing.T) {
	tc := testutil.SystemTest(t)

	datasetName := "my_new_dataset"

	b := bytes.Buffer{}

	ctx := context.Background()

	var client *bigquery.Client

	// Creates a client.
	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	//Creates dataset.
	if err := client.Dataset(datasetName).Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
		log.Fatalf("Failed to create dataset: %v", err)
	}

	defer func() {
		if err := client.Dataset(datasetName).Delete(ctx); err != nil {
			log.Fatalf("Failed to delete dataset: %v", err)
		}
	}()

	if err := viewDatasetAccessPolicies(&b, tc.ProjectID, datasetName); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Role"; !strings.Contains(got, want) {
		t.Errorf("viewDatasetAccessPolicies: expected %q to contain %q", got, want)
	}

}
