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
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/iam"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestGrantAccessView(t *testing.T) {
	tc := testutil.SystemTest(t)

	preffix := getPrefix()
	topic := fmt.Sprintf("%s_grant_access_to_view", preffix)

	datasetName := fmt.Sprintf("%s_dataset", topic)
	tableName := fmt.Sprintf("%s_table", topic)
	viewName := fmt.Sprintf("%s_view", topic)

	b := bytes.Buffer{}

	ctx := context.Background()

	// Creates Big Query client.
	client, err := testClient(t)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}

	// Creates dataset handler.
	dataset := client.Dataset(datasetName)

	// Once test is run, resources and clients are closed
	defer testCleanup(t, client, datasetName)

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
	if err := dataset.Table(viewName).Create(ctx, viewMetadata); err != nil {
		t.Errorf("Failed to create view: %v", err)
	}

	if err := grantAccessToResource(&b, tc.ProjectID, datasetName, viewName); err != nil {
		t.Error(err)
	}

	if got, want := b.String(), fmt.Sprintf("Details for Access entries in table or view %v.", viewName); !strings.Contains(got, want) {
		t.Errorf("viewTableAccessPolicies: expected %q to contain %q", got, want)
	}
}

func TestGrantAccessTable(t *testing.T) {
	tc := testutil.SystemTest(t)

	preffix := getPrefix()
	topic := fmt.Sprintf("%s_grant_access_to_table", preffix)

	datasetName := fmt.Sprintf("%s_dataset", topic)
	tableName := fmt.Sprintf("%s_table", topic)

	b := bytes.Buffer{}

	ctx := context.Background()

	// Creates Big Query client.
	client, err := testClient(t)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}

	// Creates dataset handler.
	dataset := client.Dataset(datasetName)

	// Once test is run, resources and clients are closed
	defer testCleanup(t, client, datasetName)

	// Creates dataset.
	if err := dataset.Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
		t.Errorf("Failed to create dataset: %v", err)
	}

	// Creates table.
	table := dataset.Table(tableName)
	if err := table.Create(ctx, &bigquery.TableMetadata{}); err != nil {
		t.Errorf("Failed to create table: %v", err)
	}

	if err := grantAccessToResource(&b, tc.ProjectID, datasetName, tableName); err != nil {
		t.Error(err)
	}

	if got, want := b.String(), fmt.Sprintf("Details for Access entries in table or view %v.", tableName); !strings.Contains(got, want) {
		t.Errorf("viewTableAccessPolicies: expected %q to contain %q", got, want)
	}
}

func TestGrantAccessDataset(t *testing.T) {
	tc := testutil.SystemTest(t)

	preffix := getPrefix()
	topic := fmt.Sprintf("%s_grant_access_to_dataset", preffix)

	datasetName := fmt.Sprintf("%s_dataset", topic)

	b := bytes.Buffer{}

	ctx := context.Background()

	// Creates Big Query client.
	client, err := testClient(t)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}

	// Creates dataset handler.
	dataset := client.Dataset(datasetName)

	// Once test is run, resources and clients are closed
	defer testCleanup(t, client, datasetName)

	// Creates dataset.
	if err := dataset.Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
		t.Errorf("Failed to create dataset: %v", err)
	}

	if err := grantAccessToDataset(&b, tc.ProjectID, datasetName); err != nil {
		t.Error(err)
	}

	if got, want := b.String(), fmt.Sprintf("Details for Access entries in dataset %v.\n", datasetName); !strings.Contains(got, want) {
		t.Errorf("viewTableAccessPolicies: expected %q to contain %q", got, want)
	}
}

func TestRevokeAccessDataset(t *testing.T) {
	tc := testutil.SystemTest(t)

	preffix := getPrefix()
	topic := fmt.Sprintf("%s_grant_access_to_dataset", preffix)

	datasetName := fmt.Sprintf("%s_dataset", topic)

	b := bytes.Buffer{}

	ctx := context.Background()

	entity := "example-analyst-group@google.com"

	// Creates Big Query client.
	client, err := testClient(t)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}

	// Creates dataset handler.
	dataset := client.Dataset(datasetName)

	// Once test is run, resources and clients are closed
	defer testCleanup(t, client, datasetName)

	// Creates dataset.
	if err := dataset.Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
		t.Errorf("Failed to create dataset: %v", err)
	}

	//Gets metadata
	meta, err := dataset.Metadata(ctx)
	if err != nil {
		t.Errorf("Failed to get metadata: %v", err)
	}

	// Appends a new access control entry to the existing access list.
	update := bigquery.DatasetMetadataToUpdate{
		Access: append(meta.Access, &bigquery.AccessEntry{
			Role:       bigquery.ReaderRole,
			EntityType: bigquery.GroupEmailEntity,
			Entity:     entity},
		),
	}

	// Leverage the ETag for the update to assert there's been no modifications to the
	// dataset since the metadata was originally read.
	if _, err := dataset.Update(ctx, update, meta.ETag); err != nil {
		t.Errorf("Failed to update metadata: %v", err)
	}

	if err = revokeAccessToDataset(&b, tc.ProjectID, datasetName, entity); err != nil {
		t.Error(err)
	}

	if got, want := b.String(), fmt.Sprintf("Details for Access entries in dataset %v.\n", datasetName); !strings.Contains(got, want) {
		t.Errorf("viewTableAccessPolicies: expected %q to contain %q", got, want)
	}
}

func TestRevokeTableAccessPolicies(t *testing.T) {
	tc := testutil.SystemTest(t)

	preffix := getPrefix()
	topic := fmt.Sprintf("%s_revoke_access_to_table", preffix)

	datasetName := fmt.Sprintf("%s_dataset", topic)
	tableName := fmt.Sprintf("%s_table", topic)

	b := bytes.Buffer{}

	ctx := context.Background()

	// Creates Big Query client.
	client, err := testClient(t)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}

	// Creates dataset handler.
	dataset := client.Dataset(datasetName)

	// Once test is run, resources and clients are closed
	defer testCleanup(t, client, datasetName)

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

	preffix := getPrefix()
	topic := fmt.Sprintf("%s_revoke_access_to_view", preffix)

	datasetName := fmt.Sprintf("%s_dataset", topic)
	tableName := fmt.Sprintf("%s_table", topic)
	viewName := fmt.Sprintf("%s_view", topic)

	b := bytes.Buffer{}

	ctx := context.Background()

	// Creates Big Query client.
	client, err := testClient(t)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}

	// Creates dataset handler.
	dataset := client.Dataset(datasetName)

	// Once test is run, resources and clients are closed
	defer testCleanup(t, client, datasetName)

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

func TestViewDatasetAccessPolicies(t *testing.T) {
	tc := testutil.SystemTest(t)

	preffix := getPrefix()
	topic := fmt.Sprintf("%s_view_access_to_dataset", preffix)

	datasetName := fmt.Sprintf("%s_dataset", topic)

	b := bytes.Buffer{}

	ctx := context.Background()

	// Creates Big Query client.
	client, err := testClient(t)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}

	// Creates dataset.
	if err := client.Dataset(datasetName).Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
		t.Errorf("Failed to create dataset: %v", err)
	}

	// Once test is run, resources and clients are closed
	defer testCleanup(t, client, datasetName)

	if err := viewDatasetAccessPolicies(&b, tc.ProjectID, datasetName); err != nil {
		t.Error(err)
	}

	if got, want := b.String(), "Details for Access entries in dataset"; !strings.Contains(got, want) {
		t.Errorf("viewDatasetAccessPolicies: expected %q to contain %q", got, want)
	}

}

func TestViewTableAccessPolicies(t *testing.T) {
	tc := testutil.SystemTest(t)

	preffix := getPrefix()
	topic := fmt.Sprintf("%s_view_access_to_table", preffix)

	datasetName := fmt.Sprintf("%s_dataset", topic)
	tableName := fmt.Sprintf("%s_table", topic)

	b := bytes.Buffer{}

	ctx := context.Background()

	// Creates Big Query client.
	client, err := testClient(t)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}

	// Creates dataset handler.
	dataset := client.Dataset(datasetName)

	// Once test is run, resources and clients are closed
	defer testCleanup(t, client, datasetName)

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

	preffix := getPrefix()
	topic := fmt.Sprintf("%s_view_access_to_view", preffix)

	datasetName := fmt.Sprintf("%s_dataset", topic)
	tableName := fmt.Sprintf("%s_table", topic)
	viewName := fmt.Sprintf("%s_view", topic)

	b := bytes.Buffer{}

	ctx := context.Background()

	// Creates Big Query client.
	client, err := testClient(t)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}

	// Creates dataset handler.
	dataset := client.Dataset(datasetName)

	// Once test is run, resources and clients are closed
	defer testCleanup(t, client, datasetName)

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

func getPrefix() string {
	return time.Now().Format("2006_01_02_15_04_05")
}

func testClient(t *testing.T) (*bigquery.Client, error) {
	t.Helper()

	ctx := context.Background()
	tc := testutil.SystemTest(t)

	// Creates a client.
	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	return client, nil
}

func testCleanup(t *testing.T, client *bigquery.Client, datasetName string) {
	t.Helper()

	ctx := context.Background()

	if err := client.Dataset(datasetName).DeleteWithContents(ctx); err != nil {
		t.Errorf("Failed to delete table: %v", err)
	}

	if err := client.Close(); err != nil {
		t.Fatalf("Failed to close Big Query client: %v", err)
	}
}
