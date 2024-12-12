// Copyright 2021 Google LLC
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
	"regexp"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	datacatalog "cloud.google.com/go/datacatalog/apiv1"
	"cloud.google.com/go/datacatalog/apiv1/datacatalogpb"
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

	ctx := context.Background()

	table, bqCleanup, err := createTestTable(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("couldn't setup BQ test resources: %v", err)
	}

	stdOut, stdErr, err := m.Run(nil, 30*time.Second, fmt.Sprintf("--project_id=%s", tc.ProjectID), fmt.Sprintf("--table=%s", table))
	if err != nil {
		t.Errorf("execution failed: %v", err)
	}

	// Look for a known substring in the output
	if !strings.Contains(string(stdOut), "Created tag: projects/") {
		t.Errorf("Did not find expected output.  Stdout: %s", string(stdOut))
	}

	// Handle errors during quickstart invocation.
	if len(stdErr) > 0 {
		t.Errorf("did not expect stderr output, got %d bytes: %s", len(stdErr), string(stdErr))
	}

	// cleanup data catalog resources
	if err := cleanupDataCatalog(ctx, string(stdOut)); err != nil {
		t.Errorf("failed cleanup: %v", err)
	}

	// cleanup the BQ resources we created to support the quickstart
	bqCleanup()
}

// createTestTable creates a BQ dataset and table, allowing the quickstart to use the
// entry for applying a tag.  It returns a function for cleaning up the created resources.
func createTestTable(ctx context.Context, projectID string) (string, func(), error) {
	bqClient, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return "", nil, err
	}

	datasetID, err := bqtestutil.UniqueBQName("datacatalogtests")
	if err != nil {
		return "", nil, err
	}
	tableID, err := bqtestutil.UniqueBQName("sampletable")
	if err != nil {
		return "", nil, err
	}

	dsMeta := &bigquery.DatasetMetadata{
		Location: "US", // See https://cloud.google.com/bigquery/docs/locations
	}
	if err := bqClient.Dataset(datasetID).Create(ctx, dsMeta); err != nil {
		return "", nil, err
	}

	// define the cleanup func.
	cleanup := func() {
		bqClient.Dataset(datasetID).DeleteWithContents(ctx)
		bqClient.Close()
	}

	tableMeta := &bigquery.TableMetadata{
		Schema: bigquery.Schema{
			{Name: "full_name", Type: bigquery.StringFieldType},
			{Name: "age", Type: bigquery.IntegerFieldType},
		},
	}

	tableRef := bqClient.Dataset(datasetID).Table(tableID)
	if err := tableRef.Create(ctx, tableMeta); err != nil {
		// failed table creation.  try to cleanup the dataset opportunistically before returning.
		cleanup()
		return "", nil, err
	}

	return fmt.Sprintf("%s.%s.%s", projectID, datasetID, tableID), cleanup, nil

}

// cleanupDataCatalog examines the output of the quickstart, and removes created resources.
func cleanupDataCatalog(ctx context.Context, stdOut string) error {
	re := regexp.MustCompile("Created tag template: (.*)\n")
	matches := re.FindStringSubmatch(stdOut)
	if len(matches) != 2 {
		return nil
	}

	client, err := datacatalog.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("datacatalog.NewClient: %w", err)
	}

	return client.DeleteTagTemplate(ctx, &datacatalogpb.DeleteTagTemplateRequest{
		Name:  matches[1],
		Force: true,
	})
}
