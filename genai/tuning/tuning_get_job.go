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

// [START googlegenaisdk_tuning_job_get]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// getTuningJob retrieves details of a tuning job, including tuned model and endpoint
func getTuningJob(w io.Writer, tuningJobName string) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// Eg. tuningJobName = "projects/123456789012/locations/us-central1/tuningJobs/123456789012345"
	tuningJob, err := client.Tunings.Get(ctx, tuningJobName, nil)
	if err != nil {
		return fmt.Errorf("failed to get tuning job: %w", err)
	}

	fmt.Fprintln(w, tuningJob.TunedModel.Model)
	fmt.Fprintln(w, tuningJob.TunedModel.Endpoint)
	fmt.Fprintln(w, tuningJob.Experiment)

	// Example response:
	// projects/123456789012/locations/us-central1/models/1234567890@1
	// projects/123456789012/locations/us-central1/endpoints/123456789012345
	// projects/123456789012/locations/us-central1/metadataStores/default/contexts/tuning-experiment-2025010112345678

	return nil
}

// [END googlegenaisdk_tuning_job_get]
