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

// Package automl contains samples for Google Cloud AutoML API.
package automl

// [START automl_batch_predict]
import (
	"context"
	"fmt"
	"io"

	automl "cloud.google.com/go/automl/apiv1"
	automlpb "google.golang.org/genproto/googleapis/cloud/automl/v1"
)

// batchPredict does a batch prediction.
func batchPredict(w io.Writer, projectID string, location string, modelID string, inputURI string, outputURI string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// modelID := "ICN123456789..."
	// inputURI := "gs://BUCKET_ID/path_to_your_input_csv_or_jsonl"
	// outputURI := "gs://BUCKET_ID/path_to_save_results/"

	ctx := context.Background()
	client, err := automl.NewPredictionClient(ctx)
	if err != nil {
		return fmt.Errorf("NewPredictionClient: %v", err)
	}
	defer client.Close()

	req := &automlpb.BatchPredictRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/models/%s", projectID, location, modelID),
		InputConfig: &automlpb.BatchPredictInputConfig{
			Source: &automlpb.BatchPredictInputConfig_GcsSource{
				GcsSource: &automlpb.GcsSource{
					InputUris: []string{inputURI},
				},
			},
		},
		OutputConfig: &automlpb.BatchPredictOutputConfig{
			Destination: &automlpb.BatchPredictOutputConfig_GcsDestination{
				GcsDestination: &automlpb.GcsDestination{
					OutputUriPrefix: outputURI,
				},
			},
		},
		Params: map[string]string{
			"score_threshold": "0.8", // [0.0-1.0] Only produce results higher than this value
		},
	}

	op, err := client.BatchPredict(ctx, req)
	if err != nil {
		return fmt.Errorf("BatchPredict: %v", err)
	}
	fmt.Fprintf(w, "Processing operation name: %q\n", op.Name())

	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %v", err)
	}

	fmt.Fprintf(w, "Batch Prediction results saved to Cloud Storage bucket.\n")
	fmt.Fprintf(w, "%v", resp)

	return nil
}

// [END automl_batch_predict]
