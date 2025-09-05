// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package batch_prediction shows how to use the GenAI SDK to batch prediction.

package batch_prediction

// [START googlegenaisdk_batchpredict_embeddings_with_gcs]
import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/genai"
)

// generateBatchEmbeddings shows how to run a batch embeddings prediction job.
func generateBatchEmbeddings(w io.Writer, outputURI string) error {
	// outputURI = "gs://your-bucket/your-prefix"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}
	modelName := "text-embedding-005"
	// See the documentation: https://googleapis.github.io/python-genai/genai.html#genai.batches.Batches.create
	job, err := client.Batches.Create(ctx,
		modelName,
		&genai.BatchJobSource{
			Format: "jsonl",
			// Source link: https://storage.cloud.google.com/cloud-samples-data/generative-ai/embeddings/embeddings_input.jsonl
			GCSURI: []string{"gs://cloud-samples-data/generative-ai/embeddings/embeddings_input.jsonl"},
		},
		&genai.CreateBatchJobConfig{
			Dest: &genai.BatchJobDestination{
				Format: "jsonl",
				GCSURI: outputURI,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create batch job: %w", err)
	}

	fmt.Fprintf(w, "Job name: %s\n", job.Name)
	fmt.Fprintf(w, "Job state: %s\n", job.State)
	// Example response:
	//  Job name: projects/%PROJECT_ID%/locations/us-central1/batchPredictionJobs/9876453210000000000
	//  Job state: JOB_STATE_PENDING

	// See the documentation: https://googleapis.github.io/python-genai/genai.html#genai.types.BatchJob
	completedStates := map[genai.JobState]bool{
		genai.JobStateSucceeded: true,
		genai.JobStateFailed:    true,
		genai.JobStateCancelled: true,
		genai.JobStatePaused:    true,
	}

	// Poll until job finishes
	for !completedStates[job.State] {
		time.Sleep(30 * time.Second)
		job, err = client.Batches.Get(ctx, job.Name, nil)
		if err != nil {
			return fmt.Errorf("failed to get batch job: %w", err)
		}
		fmt.Fprintf(w, "Job state: %s\n", job.State)

		if job.State == genai.JobStateFailed {
			fmt.Fprintf(w, "Error: %+v\n", job.Error)
			break
		}
	}

	//  Example response:
	//    Job state: JOB_STATE_PENDING
	//    Job state: JOB_STATE_RUNNING
	//    Job state: JOB_STATE_RUNNING
	//    ...
	//    Job state: JOB_STATE_SUCCEEDED

	return nil
}

// [END googlegenaisdk_batchpredict_embeddings_with_gcs]
