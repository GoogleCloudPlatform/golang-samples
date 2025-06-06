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

func TestGrantAccessView(t *testing.T) {
	tc := testutil.SystemTest(t)

	prefixer := testPrefix()
	prefix := fmt.Sprintf("%s_grant_access_to_view", prefixer)

	datasetName := fmt.Sprintf("%s_dataset", prefix)
	tableName := fmt.Sprintf("%s_table", prefix)
	viewName := fmt.Sprintf("%s_view", prefix)

	ctx := context.Background()

	var buf bytes.Buffer

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

	// Table schema.
	sampleSchema := bigquery.Schema{
		{Name: "id", Type: bigquery.IntegerFieldType, Required: true},
	}

	tableMetaData := &bigquery.TableMetadata{
		Schema: sampleSchema,
	}

	// Create table.
	table := dataset.Table(tableName)
	if err := table.Create(ctx, tableMetaData); err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Sets view query.
	viewMetadata := &bigquery.TableMetadata{
		ViewQuery: fmt.Sprintf("SELECT * FROM `%s.%s`", datasetName, tableName),
	}

	// Create view.
	if err := dataset.Table(viewName).Create(ctx, viewMetadata); err != nil {
		t.Fatalf("Failed to create view: %v", err)
	}

	if err := grantAccessToResource(&buf, tc.ProjectID, datasetName, viewName); err != nil {
		t.Error(err)
	}

	if got, want := buf.String(), fmt.Sprintf("Details for Access entries in table or view %v.", viewName); !strings.Contains(got, want) {
		t.Errorf("viewTableAccessPolicies: expected %q to contain %q", got, want)
	}
}

func TestGrantAccessTable(t *testing.T) {
	tc := testutil.SystemTest(t)

	prefixer := testPrefix()
	prefix := fmt.Sprintf("%s_grant_access_to_table", prefixer)

	datasetName := fmt.Sprintf("%s_dataset", prefix)
	tableName := fmt.Sprintf("%s_table", prefix)

	ctx := context.Background()

	var buf bytes.Buffer

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

	// Create table.
	table := dataset.Table(tableName)
	if err := table.Create(ctx, &bigquery.TableMetadata{}); err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	if err := grantAccessToResource(&buf, tc.ProjectID, datasetName, tableName); err != nil {
		t.Error(err)
	}

	if got, want := buf.String(), fmt.Sprintf("Details for Access entries in table or view %v.", tableName); !strings.Contains(got, want) {
		t.Errorf("viewTableAccessPolicies: expected %q to contain %q", got, want)
	}
}
