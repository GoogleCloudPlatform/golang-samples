// Copyright 2020 Google LLC
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
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMain(t *testing.T) {
	tc := testutil.SystemTest(t)
	env := map[string]string{"GOOGLE_CLOUD_PROJECT": tc.ProjectID}
	scope := fmt.Sprintf("projects/%s", tc.ProjectID)
	query := fmt.Sprintf("name:%s", tc.ProjectID)

	ctx := context.Background()
	// Creates a bigquery client.
	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("failed to create bigquery client: %v", err)
	}
	datasetID := strings.Replace(fmt.Sprintf("%s-for-asset-search", tc.ProjectID), "-", "_", -1)

	m := testutil.BuildMain(t)
	defer m.Cleanup()

	if !m.Built() {
		t.Errorf("failed to build app")
	}

	createDataset(ctx, t, client, datasetID)
	stdOut, stdErr, err := m.Run(env, 2*time.Minute, fmt.Sprintf("--scope=%s", scope), fmt.Sprintf("--query=%s", query))
	deleteDataset(ctx, t, client, datasetID)

	if err != nil {
		t.Errorf("execution failed: %v", err)
	}
	if len(stdErr) > 0 {
		t.Errorf("did not expect stderr output, got %d bytes: %s", len(stdErr), string(stdErr))
	}
	got := string(stdOut)
	if !strings.Contains(got, tc.ProjectID) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, tc.ProjectID)
	}
}

func createDataset(ctx context.Context, t *testing.T, client *bigquery.Client, datasetID string) {
	d := client.Dataset(datasetID)
	if _, err := d.Metadata(ctx); err == nil {
		deleteDataset(ctx, t, client, datasetID)
	}
	meta := &bigquery.DatasetMetadata{
		Location: "US", // See https://cloud.google.com/bigquery/docs/locations
	}
	if err := client.Dataset(datasetID).Create(ctx, meta); err != nil {
		t.Fatalf("Dataset.Create(%q): %v", datasetID, err)
	}
}

func deleteDataset(ctx context.Context, t *testing.T, client *bigquery.Client, datasetID string) {
	d := client.Dataset(datasetID)
	if errDelete := d.DeleteWithContents(ctx); errDelete != nil {
		t.Fatalf("Dataset.Delete(%q): %v", datasetID, errDelete)
	}
}
