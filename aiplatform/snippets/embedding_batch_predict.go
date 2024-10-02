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

package snippets

// [START generativeaionvertexai_sdk_embedding_batch]
import (
	"context"
	"fmt"
	"io"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/api/option"
)

// embedBatchPredict generates embeddings from text using batch processing.
func embedBatchPredict(w io.Writer, projectID, location, name, outputURI string, inputURIs []string) error {
	// inputURI := []string{"gs://cloud-samples-data/generative-ai/embeddings/embeddings_input.jsonl"}
	// outputURI: existing template path. Following formats are allowed:
	// 	- gs://BUCKET_NAME/DIRECTORY/
	// 	- bq://project_name.llm_dataset

	ctx := context.Background()
	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	// Pretrained model
	model := "publishers/google/models/textembedding-gecko"

	client, err := aiplatform.NewJobClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		return err
	}

	req := &aiplatformpb.CreateBatchPredictionJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		BatchPredictionJob: &aiplatformpb.BatchPredictionJob{
			DisplayName: name,
			Model:       model,
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

// [END generativeaionvertexai_sdk_embedding_batch]
