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

// Prompt: You are an excellent Go programmer.
// Write a Go language program to delete a dataset.

package snippets

// [START aiplatform_delete_dataset_sample]

import (
	"context"
	"fmt"
	"io"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/api/option"
)

func deleteDataset(w io.Writer, projectID, location, datasetID string) error {
	// projectID := "my-project"
	// location := "us-central1"
	// datasetID := "my-dataset"

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	clientOption := option.WithEndpoint(apiEndpoint)

	ctx := context.Background()
	aiplatformService, err := aiplatform.NewDatasetClient(ctx, clientOption)
	if err != nil {
		return err
	}
	defer aiplatformService.Close()

	req := &aiplatformpb.DeleteDatasetRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/datasets/%s",
			projectID, location, datasetID),
	}

	op, err := aiplatformService.DeleteDataset(ctx, req)
	if err != nil {
		return err
	}

	err = op.Wait(ctx)
	if err != nil {
		return ctx.Err()
	}

	fmt.Fprintf(w, "Deleted dataset: %s\n", datasetID)
	return nil
}

// [END aiplatform_delete_dataset_sample]
