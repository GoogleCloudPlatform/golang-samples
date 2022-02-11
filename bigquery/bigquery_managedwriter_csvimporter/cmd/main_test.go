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

package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/bigquery/snippets/bqtestutil"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestApp(t *testing.T) {
	tc := testutil.SystemTest(t)
	m := testutil.BuildMain(t)
	defer m.Cleanup()

	if !m.Built() {
		t.Errorf("failed to build app")
	}

	// Setup destination BQ resources
	dataset, table, cleanup, err := setupBigQueryResources(context.Background(), tc.ProjectID)
	if err != nil {
		t.Errorf("setupBigQueryResources: %v", err)
	}

	defer cleanup()
	stdOut, stdErr, err := m.Run(nil, 30*time.Second,
		fmt.Sprintf("--project_id=%s", tc.ProjectID),
		fmt.Sprintf("--dataset=%s", dataset),
		fmt.Sprintf("--table=%s", table),
		fmt.Sprintf("--max_rows=100000"),
		fmt.Sprintf("--verbose=false"))
	if err != nil {
		t.Errorf("execution failed: %v", err)
	}

	// We don't look for specific strings, just expect at least 1kb of output.
	testString := "committed data successfully"
	if !strings.Contains(string(stdOut), testString) {
		t.Errorf("expected commit message.  Stdout: %s", string(stdOut))
	}

	if len(stdErr) > 0 {
		t.Errorf("did not expect stderr output, got %d bytes: %s", len(stdErr), string(stdErr))
	}
}

func setupBigQueryResources(ctx context.Context, projectID string) (datasetID string, tableID string, cleanupF func(), err error) {
	// until we have resources, cleanup needs nothing
	cleanupF = func() {}

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return "", "", cleanupF, fmt.Errorf("failed to create bigquery client(%q): %v", projectID, err)
	}

	meta := &bigquery.DatasetMetadata{
		Location: "US",
	}
	datasetID, err = bqtestutil.UniqueBQName("managedwriter_csv_test")
	if err != nil {
		return "", "", cleanupF, fmt.Errorf("failed to generate dataset name: %v", err)
	}
	if err := client.Dataset(datasetID).Create(ctx, meta); err != nil {
		return "", "", cleanupF, fmt.Errorf("failed to create dataset %q: %v", datasetID, err)
	}
	cleanupF = func() {
		client.Dataset(datasetID).DeleteWithContents(ctx)
		client.Close()
	}

	tableID, err = bqtestutil.UniqueBQName("testtable")
	if err != nil {
		return "", "", cleanupF, fmt.Errorf("failed to generate table name: %v", err)
	}

	tableMeta := &bigquery.TableMetadata{
		Schema: bigquery.Schema{
			{Name: "date", Type: bigquery.StringFieldType},
			{Name: "airline", Type: bigquery.StringFieldType},
			{Name: "airline_code", Type: bigquery.StringFieldType},
			{Name: "departure_airport", Type: bigquery.StringFieldType},
			{Name: "departure_state", Type: bigquery.StringFieldType},
			{Name: "departure_lat", Type: bigquery.StringFieldType},
			{Name: "departure_lon", Type: bigquery.StringFieldType},
			{Name: "arrival_airport", Type: bigquery.StringFieldType},
			{Name: "arrival_state", Type: bigquery.StringFieldType},
			{Name: "arrival_lat", Type: bigquery.StringFieldType},
			{Name: "arrival_lon", Type: bigquery.StringFieldType},
			{Name: "departure_schedule", Type: bigquery.StringFieldType},
			{Name: "departure_actual", Type: bigquery.StringFieldType},
			{Name: "departure_delay", Type: bigquery.StringFieldType},
			{Name: "arrival_schedule", Type: bigquery.StringFieldType},
			{Name: "arrival_actual", Type: bigquery.StringFieldType},
			{Name: "arrival_delay", Type: bigquery.StringFieldType},
		},
	}

	if err := client.Dataset(datasetID).Table(tableID).Create(ctx, tableMeta); err != nil {
		return "", "", cleanupF, fmt.Errorf("failed to create test table(%q %q): %v", datasetID, tableID, err)
	}
	return datasetID, tableID, cleanupF, nil
}
