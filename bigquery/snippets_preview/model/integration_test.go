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

// Package model demonstrates interactions with BigQuery ML models.
package model

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"cloud.google.com/go/bigquery/v2/apiv2/bigquerypb"
	"cloud.google.com/go/bigquery/v2/apiv2_client"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const testTimeout = 30 * time.Second
const testLocation = "us-west1"

func TestModelSnippets(t *testing.T) {
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
			dsID := fmt.Sprintf("snippettesting_models_%s_%d", name, time.Now().UnixNano())

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

			modelID := fmt.Sprintf("modelsnippet_%s_%d", name, time.Now().UnixNano())

			if err := createModel(ctx, client, projID, dsID, modelID); err != nil {
				t.Fatalf("createModel failure: %v", err)
			}

			if err := updateModel(client, io.Discard, projID, dsID, modelID); err != nil {
				t.Fatalf("updateModel(%q,%q,%q): %v", projID, dsID, modelID, err)
			}
			if err := listModels(client, io.Discard, projID, dsID); err != nil {
				t.Fatalf("listModels(%q,%q): %v", projID, dsID, err)
			}
			if err := deleteModel(client, projID, dsID, modelID); err != nil {
				t.Fatalf("deleteModel(%q,%q,%q): %v", projID, dsID, modelID, err)
			}
		})
	}
}

func createModel(ctx context.Context, client *apiv2_client.Client, projectID, datasetID, modelID string) error {
	sqlID := fmt.Sprintf("`%s`.%s.%s", projectID, datasetID, modelID)
	sql := fmt.Sprintf(`
	  CREATE MODEL %s
	  OPTIONS (
	    model_type = 'linear_reg',
		max_iterations = 1,
		learn_rate=0.4,
		learn_rate_strategy='constant'
	  ) AS (
	    SELECT 'a' as f1, 2.0 as label,
		UNION ALL
		SELECT 'b' as f1, 3.8 as label
	  )`, sqlID)

	qReq := &bigquerypb.PostQueryRequest{
		ProjectId: projectID,
		QueryRequest: &bigquerypb.QueryRequest{
			Query: sql,
			UseLegacySql: &wrapperspb.BoolValue{
				Value: false,
			},
		},
	}
	resp, err := client.Query(ctx, qReq)
	if err != nil {
		return fmt.Errorf("Query: %w", err)
	}
	job := resp.GetJobReference()
	var jobDone bool
	if jc := resp.GetJobComplete(); jc != nil {
		jobDone = jc.GetValue()
	}
	for {
		if jobDone {
			return nil
		}
		pollReq := &bigquerypb.GetQueryResultsRequest{
			ProjectId: projectID,
			JobId:     job.GetJobId(),
		}
		if loc := job.GetLocation(); loc != nil {
			pollReq.Location = loc.GetValue()
		}
		resp, err := client.GetQueryResults(ctx, pollReq)
		if err != nil {
			return fmt.Errorf("GetQueryResults: %w", err)
		}
		if jc := resp.GetJobComplete(); jc != nil {
			jobDone = jc.GetValue()
		}
	}
}
