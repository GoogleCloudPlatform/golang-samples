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

// [START automl_language_entity_extraction_create_model]
import (
	"context"
	"fmt"
	"io"

	automl "cloud.google.com/go/automl/apiv1"
	automlpb "google.golang.org/genproto/googleapis/cloud/automl/v1"
)

// languageTextClassificationCreateModel creates a model for text classification.
func languageTextClassificationCreateModel(w io.Writer, projectID string, location string, datasetID string, modelName string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// datasetID := "TCN123456789..."
	// modelName := "model_display_name"

	ctx := context.Background()
	client, err := automl.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	req := &automlpb.CreateModelRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Model: &automlpb.Model{
			DisplayName: modelName,
			DatasetId:   datasetID,
			ModelMetadata: &automlpb.Model_TextClassificationModelMetadata{
				TextClassificationModelMetadata: &automlpb.TextClassificationModelMetadata{},
			},
		},
	}

	op, err := client.CreateModel(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateModel: %v", err)
	}
	fmt.Fprintf(w, "Processing operation name: %q\n", op.Name())
	fmt.Fprintf(w, "Training started...\n")

	return nil
}

// [END automl_language_entity_extraction_create_dataset]
