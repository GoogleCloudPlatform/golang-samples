// Copyright 2022 Google LLC
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

// Package client demonstrates how to setup a BigQuery client.
package client

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/bigquery/snippets/bqtestutil"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestClients(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	location := "us-east4"
	regionalEndpoint := fmt.Sprintf("https://%s-bigquery.googleapis.com/bigquery/v2/", location)
	client, err := setClientEndpoint(regionalEndpoint, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	datasetID, err := bqtestutil.UniqueBQName("golang_snippet_test_dataset")
	if err != nil {
		t.Fatal(err)
	}

	anotherLocation := "us-central1"
	dataset := client.Dataset(datasetID)
	err = dataset.Create(ctx, &bigquery.DatasetMetadata{
		Location: anotherLocation,
	})
	if err == nil {
		t.Fatalf("should fail to create dataset on location %s when pointing to location %s", anotherLocation, location)
		defer dataset.DeleteWithContents(ctx)
	}

	err = dataset.Create(ctx, &bigquery.DatasetMetadata{
		Location: location,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer dataset.DeleteWithContents(ctx)
}
