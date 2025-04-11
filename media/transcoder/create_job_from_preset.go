// Copyright 2020 Google LLC
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

package transcoder

// [START transcoder_create_job_from_preset]
import (
	"context"
	"fmt"
	"io"

	transcoder "cloud.google.com/go/video/transcoder/apiv1"
	"cloud.google.com/go/video/transcoder/apiv1/transcoderpb"
)

// createJobFromPreset creates a job based on a given preset template. See
// https://cloud.google.com/transcoder/docs/how-to/jobs#create_jobs_presets
// for more information.
func createJobFromPreset(w io.Writer, projectID string, location string, inputURI string, outputURI string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// inputURI := "gs://my-bucket/my-video-file"
	// outputURI := "gs://my-bucket/my-output-folder/"
	preset := "preset/web-hd"
	ctx := context.Background()
	client, err := transcoder.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &transcoderpb.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Job: &transcoderpb.Job{
			InputUri:  inputURI,
			OutputUri: outputURI,
			JobConfig: &transcoderpb.Job_TemplateId{
				TemplateId: preset,
			},
		},
	}
	// Creates the job, Jobs take a variable amount of time to run.
	// You can query for the job state.
	response, err := client.CreateJob(ctx, req)
	if err != nil {
		return fmt.Errorf("createJobFromPreset: %w", err)
	}

	fmt.Fprintf(w, "Job: %v", response.GetName())
	return nil
}

// [END transcoder_create_job_from_preset]
