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
	"context"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// testPrefix creates a prefix that serves as an identifier when creating a resource
// in order to isolate possible conflicts in concurrent tests.
func testPrefix() string {
	return time.Now().Format("2006_01_02_15_04_05")
}

// testClient creates a BigQuery client in order to avoid repeatedly creating customers in the common way
func testClient(t *testing.T) (*bigquery.Client, error) {
	t.Helper()

	ctx := context.Background()
	tc := testutil.SystemTest(t)

	// Create a client.
	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	return client, nil
}

// testCleanup is used for cleaning the dataset and its content once test is done
func testCleanup(t *testing.T, client *bigquery.Client, datasetName string) {
	t.Helper()

	ctx := context.Background()

	if err := client.Dataset(datasetName).DeleteWithContents(ctx); err != nil {
		t.Errorf("Failed to delete table: %v", err)
	}
}
