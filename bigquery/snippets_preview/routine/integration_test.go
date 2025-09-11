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

// Package routine provides some basic snippet examples for working with routines using
// the preview BigQuery Cloud Client Library.
package routine

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

func TestRoutineSnippet(t *testing.T) {
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
			dsID := fmt.Sprintf("snippettesting_routines_%s_%d", name, time.Now().UnixNano())

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

			routineID := fmt.Sprintf("routinesnippet_%s_%d", name, time.Now().UnixNano())

			if err := createRoutine(client, io.Discard, projID, dsID, routineID); err != nil {
				t.Fatalf("createRoutine(%q,%q,%q): %v", projID, dsID, routineID, err)
			}
			if err := getRoutine(client, io.Discard, projID, dsID, routineID); err != nil {
				t.Fatalf("getRoutine(%q,%q,%q): %v", projID, dsID, routineID, err)
			}
			if err := updateRoutine(client, io.Discard, projID, dsID, routineID); err != nil {
				t.Fatalf("updateRoutine(%q,%q,%q): %v", projID, dsID, routineID, err)
			}
			if err := listRoutines(client, io.Discard, projID, dsID); err != nil {
				t.Fatalf("listRoutines(%q,%q): %v", projID, dsID, err)
			}
			if err := deleteRoutine(client, projID, dsID, routineID); err != nil {
				t.Fatalf("deleteRoutine(%q,%q,%q): %v", projID, dsID, routineID, err)
			}
		})
	}
}
