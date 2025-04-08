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

package snippets

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/option"
)

var (
	projectID string
	location  string = "us-central1"
)

func setupDeleteDataset(t *testing.T) (datasetID string) {
	t.Helper()

	// Create a new dataset
	ctx := context.Background()
	tc := testutil.SystemTest(t)

	projectID = tc.ProjectID
	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	clientOption := option.WithEndpoint(apiEndpoint)
	client, err := aiplatform.NewDatasetClient(ctx, clientOption)
	if err != nil {
		log.Fatalf("aiplatform.NewDatasetClient: %v", err)
	}
	defer client.Close()

	op, err := client.CreateDataset(ctx, &aiplatformpb.CreateDatasetRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Dataset: &aiplatformpb.Dataset{
			DisplayName:       "my-dataset",
			MetadataSchemaUri: "gs://google-cloud-aiplatform/schema/dataset/metadata/image_1.0.0.yaml",
		},
	})

	if err != nil {
		log.Fatalf("CreateDataset() failed: %v", err)
	}

	dataset, err := op.Wait(ctx)
	if err != nil {
		log.Fatalf("Wait() failed: %v", err)
	}

	datasetName := dataset.Name
	datasetID = strings.Split(datasetName, "/")[5]
	return datasetID
}

func TestDeleteDataset(t *testing.T) {
	datasetID := setupDeleteDataset(t)

	var buf bytes.Buffer
	err := deleteDataset(&buf, projectID, location, datasetID)
	if err != nil {
		t.Fatalf("DeleteDataset() failed: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "Deleted dataset") {
		t.Errorf(
			"DeleteDataset() got %q, want to contain %q",
			got,
			"Deleted dataset",
		)
	}
}
