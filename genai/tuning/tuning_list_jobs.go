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

// Package tuning shows how to use the GenAI SDK to list tuning jobs.
package tuning

// [START googlegenaisdk_tuning_job_list]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// listTuningJobs demonstrates how to list tuning jobs.
func listTuningJobs(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	page, errList := client.Tunings.List(ctx, nil)
	if errList != nil {
		return fmt.Errorf("failed to list tuning jobs: %w", errList)
	}
	for _, job := range page.Items {
		fmt.Fprintln(w, job.Name)
	}

	// Example response:
	// projects/123456789012/locations/us-central1/tuningJobs/123456789012345
	// ...
	return nil
}

// [END googlegenaisdk_tuning_job_list]
