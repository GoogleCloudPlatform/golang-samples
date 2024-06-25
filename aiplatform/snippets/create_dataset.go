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

// [START aiplatform_create_dataset]
import (
	"context"
	"fmt"
	"io"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/api/option"
)

// createDataset creates a dataset in Vertex AI
func createDataset(w io.Writer, projectID, location string) error {
	// projectID := "my-project"
	// location := "us-central1"

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	clientOption := option.WithEndpoint(apiEndpoint)

	ctx := context.Background()
	client, err := aiplatform.NewDatasetClient(ctx, clientOption)
	if err != nil {
		return fmt.Errorf("aiplatform.NewDatasetClient: %w", err)
	}
	defer client.Close()

	// Create a new, empty image dataset
	// Vertex AI automatically assigns an ID for the dataset resource
	req := &aiplatformpb.CreateDatasetRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Dataset: &aiplatformpb.Dataset{
			DisplayName:       "my-image-dataset",
			MetadataSchemaUri: "gs://google-cloud-aiplatform/schema/dataset/metadata/image_1.0.0.yaml",
		},
	}

	resp, err := client.CreateDataset(ctx, req)
	if err != nil {
		return err
	}

	// Wait for the longrunning operation to complete
	dataset, err := resp.Wait(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Created dataset with resource name:")
	fmt.Fprintf(w, "%s\n", dataset.Name)
	return nil
}

// [END aiplatform_create_dataset]
