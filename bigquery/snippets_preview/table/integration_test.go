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

// Package table provides some basic snippet examples for working with tables using
// the preview BigQuery Cloud Client Library.
package table

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"cloud.google.com/go/bigquery/v2/apiv2/bigquerypb"
	"cloud.google.com/go/bigquery/v2/apiv2_client"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const testTimeout = 30 * time.Second
const testLocation = "us-west1"

func TestTableSnippet(t *testing.T) {
	tc := testutil.SystemTest(t)
	names := []string{"gRPC", "REST"}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
			defer cancel()
			// Setup client.
			var client *apiv2_client.Client
			var err error
			if name == "gRPC" {
				client, err = apiv2_client.NewClient(ctx)
			} else {
				client, err = apiv2_client.NewRESTClient(ctx)
			}
			if err != nil {
				t.Fatalf("client creation failed: %v", err)
			}
			defer client.Close()

			// Create a test dataset.
			projID := tc.ProjectID
			dsID := fmt.Sprintf("snippettesting_tables_%s_%d", name, time.Now().UnixNano())

			_, err = client.InsertDataset(ctx, &bigquerypb.InsertDatasetRequest{
				ProjectId: projID,
				Dataset: &bigquerypb.Dataset{
					DatasetReference: &bigquerypb.DatasetReference{
						ProjectId: projID,
						DatasetId: dsID,
					},
					Location: testLocation,
				},
			})
			if err != nil {
				t.Fatalf("couldn't create test dataset: %v", err)
			}
			defer client.DeleteDataset(ctx, &bigquerypb.DeleteDatasetRequest{
				ProjectId:      projID,
				DatasetId:      dsID,
				DeleteContents: true,
			})

			tableID := fmt.Sprintf("tablesnippet_%s_%d", name, time.Now().UnixNano())

			if err := createTable(client, io.Discard, projID, dsID, tableID); err != nil {
				t.Fatalf("createTable(%q,%q,%q): %v", projID, dsID, tableID, err)
			}
			if err := updateTable(client, io.Discard, projID, dsID, tableID); err != nil {
				t.Fatalf("updateTable(%q,%q,%q): %v", projID, dsID, tableID, err)
			}
			if err := listTables(client, io.Discard, projID, dsID); err != nil {
				t.Fatalf("listTables(%q,%q): %v", projID, dsID, err)
			}
			if err := deleteTable(client, projID, dsID, tableID); err != nil {
				t.Fatalf("deleteTable(%q,%q,%q): %v", projID, dsID, tableID, err)
			}
		})
	}
}
