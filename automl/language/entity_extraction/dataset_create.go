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

// Package automl contains samples for Google Cloud AutoML API v1.
package automl

// [START automl_language_entity_extraction_create_dataset]
import (
	"context"
	"fmt"
	"io"

	automl "cloud.google.com/go/automl/apiv1"
	automlpb "google.golang.org/genproto/googleapis/cloud/automl/v1"
)

// languageEntityExtractionCreateDataset creates a dataset for text entity extraction.
func languageEntityExtractionCreateDataset(w io.Writer, projectID string, location string, datasetName string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// datasetName := "dataset_display_name"

	ctx := context.Background()
	client, err := automl.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	req := &automlpb.CreateDatasetRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Dataset: &automlpb.Dataset{
			DisplayName: datasetName,
			DatasetMetadata: &automlpb.Dataset_TextExtractionDatasetMetadata{
				TextExtractionDatasetMetadata: &automlpb.TextExtractionDatasetMetadata{},
			},
		},
	}

	op, err := client.CreateDataset(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateDataset: %v", err)
	}
	fmt.Fprintf(w, "Processing operation name: %q\n", op.Name())

	dataset, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %v", err)
	}

	fmt.Fprintf(w, "Dataset name: %v\n", dataset.GetName())

	return nil
}

// [END automl_language_entity_extraction_create_dataset]
