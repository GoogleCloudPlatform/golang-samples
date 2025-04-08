// Copyright 2019 Google LLC
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

// Package models demonstrates interactions with BigQuery's ML model
// functionality.  The models API in BigQuery does not allow direct creation
// of models,  which are instead created via CREATE MODEL queries.
package models

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/bigquery/snippets/bqtestutil"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestModels(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	datasetID, err := bqtestutil.UniqueBQName("golang_example_dataset_model")
	if err != nil {
		t.Errorf("UniqueBQName:%v", err)
	}
	if err := client.Dataset(datasetID).Create(ctx,
		&bigquery.DatasetMetadata{
			Location: "US",
		}); err != nil {
		t.Errorf("dataset.Create(%q): %v", datasetID, err)
	}
	defer client.Dataset(datasetID).DeleteWithContents(ctx)

	modelID, err := bqtestutil.UniqueBQName("golang_example_model")
	if err != nil {
		t.Errorf("UniqueBQName: %v", err)
	}
	modelRef := fmt.Sprintf("%s.%s.%s", tc.ProjectID, datasetID, modelID)

	// Create a ML model via a query.
	sql := fmt.Sprintf(`
	CREATE MODEL `+"`%s`"+`
	OPTIONS (
		model_type='linear_reg',
		max_iteration=1,
		learn_rate=0.4,
		learn_rate_strategy='constant'
	) AS (
		SELECT 'a' AS f1, 2.0 AS label
		UNION ALL
		SELECT 'b' AS f1, 3.8 AS label
	)`, modelRef)
	job, err := client.Query(sql).Run(ctx)
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}
	_, err = job.Wait(ctx)
	if err != nil {
		t.Fatalf("waiting for job completion failed: %v", err)
	}
	if err := printModelInfo(ioutil.Discard, tc.ProjectID, datasetID, modelID); err != nil {
		t.Errorf("printModelInfo(%q %q): %v", datasetID, modelID, err)
	}
	if err := listModels(ioutil.Discard, tc.ProjectID, datasetID); err != nil {
		t.Errorf("listModels(%q): %v", datasetID, err)
	}
	if err := updateModelDescription(tc.ProjectID, datasetID, modelID); err != nil {
		t.Errorf("updateModelDescription(%q %q): %v", datasetID, modelID, err)
	}

	if err := deleteModel(tc.ProjectID, datasetID, modelID); err != nil {
		t.Errorf("deleteModel(%q %q): %v", datasetID, modelID, err)
	}
}
