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

package viewiampolicytableorview

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
	testfunctions "github.com/GoogleCloudPlatform/golang-samples/bigquery/bigquery_access_quickstart"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestViewTableAccessPolicies(t *testing.T) {
	tc := testutil.SystemTest(t)

	datasetName := "my_new_dataset_go"
	tableName := "my_table"

	b := bytes.Buffer{}

	ctx := context.Background()

	// Creates bq client.
	client, err := testfunctions.TestClient(t, ctx)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}

	// Creates dataset handler.
	dataset := client.Dataset(datasetName)

	// Once test is run, resources and clients are closed
	defer testfunctions.TestCleanup(t, ctx, client, datasetName)

	//Creates dataset.
	if err := dataset.Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
		t.Errorf("Failed to create dataset: %v", err)
	}

	//Creates table.
	if err := dataset.Table(tableName).Create(ctx, &bigquery.TableMetadata{}); err != nil {
		t.Errorf("Failed to create table: %v", err)
	}

	if err := viewTableOrViewccessPolicies(&b, tc.ProjectID, datasetName, tableName); err != nil {
		t.Error(err)
	}

	if got, want := b.String(), fmt.Sprintf("Details for Access entries in table or view %v.", tableName); !strings.Contains(got, want) {
		t.Errorf("viewTableAccessPolicies: expected %q to contain %q", got, want)
	}
}

func TestViewViewAccessPolicies(t *testing.T) {
	tc := testutil.SystemTest(t)

	datasetName := "my_new_dataset_go"
	tableName := "my_table"
	viewName := "my_view"

	b := bytes.Buffer{}

	ctx := context.Background()

	// Creates bq client.
	client, err := testfunctions.TestClient(t, ctx)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}

	// Creates dataset handler.
	dataset := client.Dataset(datasetName)

	// Once test is run, resources and clients are closed
	defer testfunctions.TestCleanup(t, ctx, client, datasetName)

	// Creates dataset.
	if err := dataset.Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
		t.Errorf("Failed to create dataset: %v", err)
	}

	// Table schema.
	sampleSchema := bigquery.Schema{
		{Name: "full_name", Type: bigquery.StringFieldType},
		{Name: "age", Type: bigquery.IntegerFieldType},
	}

	tableMetaData := &bigquery.TableMetadata{
		Schema: sampleSchema,
	}

	// Creates table.
	table := dataset.Table(tableName)
	if err := table.Create(ctx, tableMetaData); err != nil {
		t.Errorf("Failed to create table: %v", err)
	}

	// Sets view query.
	viewMetadata := &bigquery.TableMetadata{
		ViewQuery: fmt.Sprintf("SELECT UPPER(full_name) FROM `%s.%s` ORDER BY age ASC", datasetName, tableName),
	}

	// Creates view
	if err := dataset.Table(viewName).Create(ctx, viewMetadata); err != nil {
		t.Errorf("Failed to create view: %v", err)
	}

	if err := viewTableOrViewccessPolicies(&b, tc.ProjectID, datasetName, viewName); err != nil {
		t.Error(err)
	}

	if got, want := b.String(), fmt.Sprintf("Details for Access entries in table or view %v.", viewName); !strings.Contains(got, want) {
		t.Errorf("viewTableAccessPolicies: expected %q to contain %q", got, want)
	}
}
