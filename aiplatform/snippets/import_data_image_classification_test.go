// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Prompt: You are a Go programmer that knows Google Cloud. Write a test for importDataImageClassification that is similar to TestCreateDataset.

package snippets

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var (
	datasetName string
	datasetID   string
	gcsURI      string
)

func setupImportDatasetImageClassification(t *testing.T) func() {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	// Setup: create dataset
	client, err := aiplatform.NewDatasetClient(ctx)
	if err != nil {
		t.Fatalf("aiplatform.NewDatasetClient: %v", err)
	}
	defer client.Close()

	tc := testutil.SystemTest(t)

	req := &aiplatformpb.CreateDatasetRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", tc.ProjectID, region),
		Dataset: &aiplatformpb.Dataset{
			DisplayName:       "my-image-dataset",
			MetadataSchemaUri: "gs://google-cloud-aiplatform/schema/dataset/metadata/image_1.0.0.yaml",
		},
	}

	op, err := client.CreateDataset(ctx, req)
	if err != nil {
		t.Fatalf("CreateDataset: %v", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		t.Fatalf("Wait: %v", err)
	}

	datasetName = resp.GetName()
	datasetID = strings.Split("/", datasetName)[5]

	// Setup: upload dataset JSONL file to GCS
	sc, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer sc.Close()

	b := testutil.CreateTestBucket(ctx, t, sc, tc.ProjectID, "vertex-image-classification")

	bucket := sc.Bucket(b)
	f, err := os.Open("../testdata/icn-dataset.jsonl")
	if err != nil {
		t.Fatalf("os.Open: %v", err)
	}
	defer f.Close()

	o := bucket.Object("vertex-image-classification/icn-dataset.jsonl")
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		t.Errorf("io.Copy: %v", err)
	}

	return func() {
		dr := &aiplatformpb.DeleteDatasetRequest{
			Name: datasetName,
		}
		op, err := client.DeleteDataset(ctx, dr)
		if err != nil {
			t.Errorf("DeleteDataset: %v", err)
		}
		if err := op.Wait(ctx); err != nil {
			t.Errorf("op.Wait: %v", err)
		}
	}
}

func TestImportDataImageClassification(t *testing.T) {
	t.Skip("skipped, see context at https://github.com/GoogleCloudPlatform/golang-samples/issues/3579")
	tc := testutil.SystemTest(t)
	teardown := setupImportDatasetImageClassification(t)
	t.Cleanup(teardown)

	var buf bytes.Buffer
	if err := importDataImageClassification(&buf, tc.ProjectID, region, datasetID, gcsURI); err != nil {
		t.Errorf("importDataImageClassification got err: %v", err)
	}
}
