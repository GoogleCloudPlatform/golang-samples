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

// Prompt: Write a code sample that creates a dataset in Vertex AI

import (
	"context"
	"fmt"
	"io"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	aiplatformpb "google.golang.org/genproto/googleapis/cloud/aiplatform/v1"
)

// createDataset creates a dataset in Vertex AI
func createDataset(w io.Writer, projectID, location, datasetID string) error {
	// projectID := "my-project"
	// location := "us-central1"
	// datasetID := "my-dataset"

	ctx := context.Background()
	client, err := aiplatform.NewDatasetClient(ctx)
	if err != nil {
		return fmt.Errorf("aiplatform.NewDatasetClient: %v", err)
	}
	defer client.Close()

	req := &aiplatformpb.CreateDatasetRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Dataset: &aiplatformpb.Dataset{
			DisplayName: datasetID,
		},
	}

	resp, err := client.CreateDataset(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateDataset: %v", err)
	}

	fmt.Fprintf(w, "Created dataset: %s\n", resp.GetName())
	return nil
}
