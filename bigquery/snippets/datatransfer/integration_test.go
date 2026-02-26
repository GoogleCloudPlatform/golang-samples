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

// Package datatransfer demonstrates interactions with the BigQuery
// Data Transfer client.
package client

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	datatransfer "cloud.google.com/go/bigquery/datatransfer/apiv1"
	"cloud.google.com/go/bigquery/datatransfer/apiv1/datatransferpb"
	"github.com/GoogleCloudPlatform/golang-samples/bigquery/snippets/bqtestutil"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	gax "github.com/googleapis/gax-go/v2"
)

func TestDataTransfer(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	dtc, err := datatransfer.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer dtc.Close()

	datasetID, err := bqtestutil.UniqueBQName("golang_snippet_test_dataset")
	if err != nil {
		t.Fatal(err)
	}

	dataset := client.Dataset(datasetID)
	if err := dataset.Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
		t.Fatal(err)
	}
	defer dataset.DeleteWithContents(ctx)

	query := `SELECT CURRENT_TIMESTAMP() as current_time, @run_time as intended_run_time,
	@run_date as intended_run_date, 17 as some_integer`
	err = createScheduledQuery(tc.ProjectID, datasetID, query)
	if err != nil {
		t.Fatal(err)
	}

	it := dtc.ListTransferConfigs(ctx, &datatransferpb.ListTransferConfigsRequest{
		Parent:        fmt.Sprintf("projects/%s", tc.ProjectID),
		DataSourceIds: []string{"scheduled_query"},
	})
	transferConfig, err := it.Next()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = deleteScheduledQuery(transferConfig.Name)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Give some time for it to run
	time.Sleep(10 * time.Second)
	err = waitTransferRun(t, dtc, transferConfig.Name)
	if err != nil {
		t.Fatal(err)
	}
}

func waitTransferRun(t *testing.T, dtc *datatransfer.Client, transferConfigID string) error {
	ctx := context.Background()
	retries := 10
	backoff := gax.Backoff{
		Initial:    1 * time.Second,
		Multiplier: 2,
		Max:        30 * time.Second,
	}
	for {
		runsIt := dtc.ListTransferRuns(ctx, &datatransferpb.ListTransferRunsRequest{
			Parent: transferConfigID,
		})

		run, err := runsIt.Next()
		if err != nil {
			t.Fatal(err)
		}
		if run.State == datatransferpb.TransferState_SUCCEEDED {
			return nil
		}
		if run.State == datatransferpb.TransferState_FAILED || run.State == datatransferpb.TransferState_CANCELLED {
			return fmt.Errorf("transfer run failed with status: %s", run.State)
		}
		retries--
		if retries <= 0 {
			break
		}
		t := backoff.Pause()
		time.Sleep(t)
	}
	return fmt.Errorf("timeout waiting for transfer run execution")
}
