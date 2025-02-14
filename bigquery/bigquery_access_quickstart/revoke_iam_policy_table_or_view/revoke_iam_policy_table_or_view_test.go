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

package revokeiampolicytableorview

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/iam"
	testfunctions "github.com/GoogleCloudPlatform/golang-samples/bigquery/bigquery_access_quickstart"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestRevokeTableAccessPolicies(t *testing.T) {
	tc := testutil.SystemTest(t)

	preffix := testfunctions.RandString()
	topic := fmt.Sprintf("%s_revoke_access_to_table", preffix)

	datasetName := fmt.Sprintf("%s_dataset", topic)
	tableName := fmt.Sprintf("%s_table", topic)

	b := bytes.Buffer{}

	ctx := context.Background()

	// Creates Big Query client.
	client, err := testfunctions.TestClient(t)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}

	// Creates dataset handler.
	dataset := client.Dataset(datasetName)

	// Once test is run, resources and clients are closed
	defer testfunctions.TestCleanup(t, client, datasetName)

	// Creates dataset.
	if err := dataset.Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
		t.Errorf("Failed to create dataset: %v", err)
	}

	// Creates table.
	table := dataset.Table(tableName)
	if err := table.Create(ctx, &bigquery.TableMetadata{}); err != nil {
		t.Errorf("Failed to create table: %v", err)
	}

	// Gets resource policy.
	policy, err := table.IAM().Policy(ctx)
	if err != nil {
		t.Errorf("Failed to get policy: %v", err)
	}

	// Adds new policy which will be deleted.
	analystEmail := "example-analyst-group@google.com"
	policy.Add(fmt.Sprintf("group:%s", analystEmail), iam.Viewer)

	// Updates resource's policy.
	err = table.IAM().SetPolicy(ctx, policy)
	if err != nil {
		t.Errorf("Failed to set policy: %v", err)
	}

	if err := revokeTableOrViewAccessPolicies(&b, tc.ProjectID, datasetName, tableName); err != nil {
		t.Error(err)
	}

	if got, want := b.String(), fmt.Sprintf("Details for Access entries in table or view %v.", tableName); !strings.Contains(got, want) {
		t.Errorf("viewTableAccessPolicies: expected %q to contain %q", got, want)
	}
}

func TestRevokeViewAccessPolicies(t *testing.T) {
	tc := testutil.SystemTest(t)

	preffix := testfunctions.RandString()
	topic := fmt.Sprintf("%s_revoke_access_to_view", preffix)

	datasetName := fmt.Sprintf("%s_dataset", topic)
	tableName := fmt.Sprintf("%s_table", topic)
	viewName := fmt.Sprintf("%s_view", topic)

	b := bytes.Buffer{}

	ctx := context.Background()

	// Creates Big Query client.
	client, err := testfunctions.TestClient(t)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}

	// Creates dataset handler.
	dataset := client.Dataset(datasetName)

	// Once test is run, resources and clients are closed
	defer testfunctions.TestCleanup(t, client, datasetName)

	// Creates dataset.
	if err := dataset.Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
		t.Errorf("Failed to create dataset: %v", err)
	}

	// Table schema.
	sampleSchema := bigquery.Schema{
		{Name: "id", Type: bigquery.IntegerFieldType, Required: true},
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
		ViewQuery: fmt.Sprintf("SELECT * FROM `%s.%s`", datasetName, tableName),
	}

	// Creates view
	view := dataset.Table(viewName)
	if err := view.Create(ctx, viewMetadata); err != nil {
		t.Errorf("Failed to create view: %v", err)
	}

	// Gets view policy.
	policy, err := table.IAM().Policy(ctx)
	if err != nil {
		t.Errorf("Failed to get policy: %v", err)
	}

	// Adds new policy which will be deleted.
	analystEmail := "example-analyst-group@google.com"
	policy.Add(fmt.Sprintf("group:%s", analystEmail), iam.Viewer)

	// Updates views's policy.
	err = table.IAM().SetPolicy(ctx, policy)
	if err != nil {
		t.Errorf("Failed to set policy: %v", err)
	}

	if err := revokeTableOrViewAccessPolicies(&b, tc.ProjectID, datasetName, viewName); err != nil {
		t.Error(err)
	}

	if got, want := b.String(), fmt.Sprintf("Details for Access entries in table or view %v.", viewName); !strings.Contains(got, want) {
		t.Errorf("viewTableAccessPolicies: expected %q to contain %q", got, want)
	}
}
