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

// createTestDataset sets up an ephemeral test dataset for use in the quickstart.  Call the returned function
// afterwards to clean up the dataset.
func createTestDataset(ctx context.Context, projectID string) (string, error, func()) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return "", err, nil
	}
	dsName, err := bqtestutil.UniqueBQName("analyticshub_quickstart_dataset")
	if err != nil {
		return "", err, nil
	}
	dataset := client.Dataset(dsName)
	if err := dataset.Create(ctx, nil); err != nil {
		return "", err, nil
	}
	return fmt.Sprintf("projects/%s/datasets/%s", projectID, dsName), nil, func() {
		dataset.DeleteWithContents(ctx)
		client.Close()
	}
}

func TestApp(t *testing.T) {
	tc := testutil.SystemTest(t)
	m := testutil.BuildMain(t)
	defer m.Cleanup()

	if !m.Built() {
		t.Errorf("failed to build app")
	}

	ctx := context.Background()
	dataset, err, cleanup := createTestDataset(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("failed to setup test dataset before running quickstart: %v", err)
	}
	defer cleanup()

	stdOut, stdErr, err := m.Run(nil, 30*time.Second,
		fmt.Sprintf("--project_id=%s", tc.ProjectID),
		fmt.Sprintf("--dataset_source=%s", dataset))
	if err != nil {
		t.Errorf("execution failed: %v", err)
	}

	// Look for a known substring in the output
	if !strings.Contains(string(stdOut), "Quickstart completed.") {
		t.Errorf("Did not find expected output.  Stdout: %s", string(stdOut))
	}

	if len(stdErr) > 0 {
		t.Errorf("did not expect stderr output, got %d bytes: %s", len(stdErr), string(stdErr))
	}
}
