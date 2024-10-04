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

// [START generativeaionvertexai_batch_text_predict]
import (
	"context"
	"fmt"
	"io"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

// batchTextPredict perform batch text prediction using a pre-trained text generation model
func batchTextPredict(w io.Writer, projectID, location, name, outputURI string, inputURIs []string) error {
	// inputURI := []string{"gs://cloud-samples-data/batch/prompt_for_batch_text_predict.jsonl"}
	// outputURI: existing template path. Following formats are allowed:
	// 	- gs://BUCKET_NAME/DIRECTORY/
	// 	- bq://project_name.llm_dataset

	ctx := context.Background()
	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	// Pretrained text model
	model := "publishers/google/models/text-bison"
	parameters := map[string]interface{}{
		"temperature":     0.2,
		"maxOutputTokens": 200,
		"topP":            0.95,
		"topK":            40,
	}

	parametersValue, err := structpb.NewValue(parameters)
	if err != nil {
		fmt.Fprintf(w, "unable to convert parameters to Value: %v", err)
		return err
	}

	client, err := aiplatform.NewJobClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		return err
	}

	req := &aiplatformpb.CreateBatchPredictionJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		BatchPredictionJob: &aiplatformpb.BatchPredictionJob{
			DisplayName:     name,
			Model:           model,
			ModelParameters: parametersValue,
			InputConfig: &aiplatformpb.BatchPredictionJob_InputConfig{
				Source: &aiplatformpb.BatchPredictionJob_InputConfig_GcsSource{
					GcsSource: &aiplatformpb.GcsSource{
						Uris: inputURIs,
					},
				},
				// List of supported formarts: https://cloud.google.com/vertex-ai/docs/reference/rpc/google.cloud.aiplatform.v1#model
				InstancesFormat: "jsonl",
			},
			OutputConfig: &aiplatformpb.BatchPredictionJob_OutputConfig{
				Destination: &aiplatformpb.BatchPredictionJob_OutputConfig_GcsDestination{
					GcsDestination: &aiplatformpb.GcsDestination{
						OutputUriPrefix: outputURI,
					},
				},
				// List of supported formarts: https://cloud.google.com/vertex-ai/docs/reference/rpc/google.cloud.aiplatform.v1#model
				PredictionsFormat: "jsonl",
			},
		},
	}

	job, err := client.CreateBatchPredictionJob(ctx, req)
	if err != nil {
		return err
	}
	fmt.Fprint(w, job.GetDisplayName())

	return nil
}

// [END generativeaionvertexai_batch_text_predict]
