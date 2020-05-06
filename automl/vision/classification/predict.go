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

// [START automl_vision_classification_predict]
import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	automl "cloud.google.com/go/automl/apiv1"
	automlpb "google.golang.org/genproto/googleapis/cloud/automl/v1"
)

// visionClassificationPredict does a prediction for image classification.
func visionClassificationPredict(w io.Writer, projectID string, location string, modelID string, filePath string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// modelID := "ICN123456789..."
	// filePath := "path/to/image.jpg"

	ctx := context.Background()
	client, err := automl.NewPredictionClient(ctx)
	if err != nil {
		return fmt.Errorf("NewPredictionClient: %v", err)
	}
	defer client.Close()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Open: %v", err)
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("ReadAll: %v", err)
	}

	req := &automlpb.PredictRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/models/%s", projectID, location, modelID),
		Payload: &automlpb.ExamplePayload{
			Payload: &automlpb.ExamplePayload_Image{
				Image: &automlpb.Image{
					Data: &automlpb.Image_ImageBytes{
						ImageBytes: bytes,
					},
				},
			},
		},
		// Params is additional domain-specific parameters.
		Params: map[string]string{
			// score_threshold is used to filter the result.
			"score_threshold": "0.8",
		},
	}

	resp, err := client.Predict(ctx, req)
	if err != nil {
		return fmt.Errorf("Predict: %v", err)
	}

	for _, payload := range resp.GetPayload() {
		fmt.Fprintf(w, "Predicted class name: %v\n", payload.GetDisplayName())
		fmt.Fprintf(w, "Predicted class score: %v\n", payload.GetClassification().GetScore())
	}

	return nil
}

// [END automl_vision_classification_predict]
