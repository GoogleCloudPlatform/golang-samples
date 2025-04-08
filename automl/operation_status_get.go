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

// [START automl_get_operation_status]
import (
	"context"
	"fmt"
	"io"

	automl "cloud.google.com/go/automl/apiv1"
	"cloud.google.com/go/automl/apiv1/automlpb"
)

// getOperationStatus gets an operation's status.
func getOperationStatus(w io.Writer, projectID string, location string, datasetID string, modelName string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// datasetID := "ICN123456789..."
	// modelName := "model_display_name"

	ctx := context.Background()
	client, err := automl.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &automlpb.CreateModelRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Model: &automlpb.Model{
			DisplayName: modelName,
			DatasetId:   datasetID,
			ModelMetadata: &automlpb.Model_ImageClassificationModelMetadata{
				ImageClassificationModelMetadata: &automlpb.ImageClassificationModelMetadata{
					TrainBudgetMilliNodeHours: 1000, // 1000 milli-node hours are 1 hour
				},
			},
		},
	}

	op, err := client.CreateModel(ctx, req)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Name: %v\n", op.Name())

	// Wait for the longrunning operation complete.
	resp, err := op.Wait(ctx)
	if err != nil && !op.Done() {
		fmt.Println("failed to fetch operation status", err)
		return err
	}
	if err != nil && op.Done() {
		fmt.Println("operation completed with error", err)
		return err
	}
	fmt.Fprintf(w, "Response: %v\n", resp)

	return nil
}

// [END automl_get_operation_status]
