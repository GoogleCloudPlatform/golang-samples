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

// [START generativeaionvertexai_batch_predict_gemini_createjob_bigquery]
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

// batchPredictBQ submits a batch prediction job using BigQuery data source as its input
func batchPredictBQ(w io.Writer, projectID, location string, inputURI string, outputURI string) error {
	// location  := "us-central1"
	// inputURI  := "bq://storage-samples.generative_ai.batch_requests_for_multimodal_input"
	// outputURI := "bq://<cloud-project-name>.<dataset-name>.<table-name>"
	modelName := "gemini-1.5-pro-002"
	jobName := "batch-predict-bq-test-001"

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
				Source: &aiplatformpb.BatchPredictionJob_InputConfig_BigquerySource{
					BigquerySource: &aiplatformpb.BigQuerySource{
						InputUri: inputURI,
					},
				},
				InstancesFormat: "bigquery",
			},

			OutputConfig: &aiplatformpb.BatchPredictionJob_OutputConfig{
				Destination: &aiplatformpb.BatchPredictionJob_OutputConfig_BigqueryDestination{
					BigqueryDestination: &aiplatformpb.BigQueryDestination{
						OutputUri: outputURI,
					},
				},
				PredictionsFormat: "bigquery",
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
			return nil
		} else {
			fmt.Fprintf(w, "batch predict job is running... job state is %s\n", job.GetState())
		}
	}
}

// [END generativeaionvertexai_batch_predict_gemini_createjob_bigquery]
