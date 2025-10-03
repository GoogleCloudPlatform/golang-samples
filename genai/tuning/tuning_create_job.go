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

// Package tuning shows how to use the GenAI SDK for tuning jobs.
package tuning

// [START googlegenaisdk_tuning_job_create]
import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/genai"
)

// createTuningJob creates a supervised tuning job using training and validation datasets.
func createTuningJob(w io.Writer, outputGCSURI string) error {
	// outputURI = "gs://your-bucket/your-prefix"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1beta1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// Training dataset (JSONL in GCS)
	training := &genai.TuningDataset{
		GCSURI: "gs://cloud-samples-data/ai-platform/generative_ai/gemini/text/sft_train_data.jsonl",
	}

	validation := &genai.TuningValidationDataset{
		GCSURI: "gs://cloud-samples-data/ai-platform/generative_ai/gemini/text/sft_validation_data.jsonl",
	}

	// Config for the tuning job
	config := &genai.CreateTuningJobConfig{
		TunedModelDisplayName: "Example tuning job",
		ValidationDataset:     validation,
	}

	// Start tuning job
	job, err := client.Tunings.Tune(ctx, "gemini-2.5-flash", training, config)
	if err != nil {
		return fmt.Errorf("failed to create tuning job: %w", err)
	}

	// Poll until the job leaves running states
	for job.State == genai.JobStateQueued || job.State == genai.JobStatePending || job.State == genai.JobStateRunning {
		time.Sleep(30 * time.Second)

		job, err = client.Tunings.Get(ctx, job.Name, nil)
		if err != nil {
			return fmt.Errorf("failed to get tuning job: %w", err)
		}
		fmt.Fprintln(w, "Job state:", job.State)
	}

	// Print results when finished
	if job.TunedModel != nil {
		fmt.Fprintln(w, "Tuned model:", job.TunedModel.Model)
		fmt.Fprintln(w, "Endpoint:", job.TunedModel.Endpoint)
	}
	fmt.Fprintln(w, "Final state:", job.State)

	// Example response:
	// Checkpoint 1:  checkpoint_id='1' epoch=1 step=10 endpoint='projects/123456789012/locations/us-central1/endpoints/123456789000000'
	// Checkpoint 2:  checkpoint_id='2' epoch=2 step=20 endpoint='projects/123456789012/locations/us-central1/endpoints/123456789012345'

	return nil
}

// [END googlegenaisdk_tuning_job_create]
