// Copyright 2024 Google LLC
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

package batchpredict

// [START generativeaionvertexai_batch_predict_gemini_createjob]
import (
	"context"
	"fmt"
	"io"
	"time"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1/aiplatformpb"

	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

// batchPredictGCS submits a batch prediction job using GCS data source as its input
func batchPredictGCS(w io.Writer, projectID, location string, inputURIs []string, outputURI string) error {
	// location := "us-central1"
	// inputURIs := []string{"gs://cloud-samples-data/batch/prompt_for_batch_gemini_predict.jsonl"}
	// outputURI := "gs://<cloud-bucket-name>/<prefix-name>"
	modelName := "gemini-1.5-pro-002"
	jobName := "batch-predict-gcs-test-001"

	ctx := context.Background()
	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	client, err := aiplatform.NewJobClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		return fmt.Errorf("unable to create aiplatform client: %w", err)
	}
	defer client.Close()

	modelParameters, err := structpb.NewValue(map[string]interface{}{
		"temperature":     0.2,
		"maxOutputTokens": 200,
	})
	if err != nil {
		return fmt.Errorf("unable to convert model parameters to protobuf value: %w", err)
	}

	req := &aiplatformpb.CreateBatchPredictionJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		BatchPredictionJob: &aiplatformpb.BatchPredictionJob{
			DisplayName:     jobName,
			Model:           fmt.Sprintf("publishers/google/models/%s", modelName),
			ModelParameters: modelParameters,
			// Check the API reference for `BatchPredictionJob` for supported input and output formats:
			// https://cloud.google.com/vertex-ai/docs/reference/rpc/google.cloud.aiplatform.v1#google.cloud.aiplatform.v1.BatchPredictionJob
			InputConfig: &aiplatformpb.BatchPredictionJob_InputConfig{
				Source: &aiplatformpb.BatchPredictionJob_InputConfig_GcsSource{
					GcsSource: &aiplatformpb.GcsSource{
						Uris: inputURIs,
					},
				},
				InstancesFormat: "jsonl",
			},
			OutputConfig: &aiplatformpb.BatchPredictionJob_OutputConfig{
				Destination: &aiplatformpb.BatchPredictionJob_OutputConfig_GcsDestination{
					GcsDestination: &aiplatformpb.GcsDestination{
						OutputUriPrefix: outputURI,
					},
				},
				PredictionsFormat: "jsonl",
			},
		},
	}

	job, err := client.CreateBatchPredictionJob(ctx, req)
	if err != nil {
		return err
	}
	fullJobId := job.GetName()
	fmt.Fprintf(w, "submitted batch predict job %q for model %q\n", fullJobId, job.GetModel())
	fmt.Fprintf(w, "job state is %s\n", job.GetState())

	for {
		time.Sleep(5 * time.Second)

		job, err := client.GetBatchPredictionJob(ctx, &aiplatformpb.GetBatchPredictionJobRequest{
			Name: fullJobId,
		})
		if err != nil {
			return fmt.Errorf("error: couldn't get updated job state: %w", err)
		}

		if job.GetEndTime() != nil {
			fmt.Fprintf(w, "batch predict job finished with state %s\n", job.GetState())
			break
		} else {
			fmt.Fprintf(w, "batch predict job is running... job state is %s\n", job.GetState())
		}
	}

	return nil
}

// [END generativeaionvertexai_batch_predict_gemini_createjob]
