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

// [START transcoder_get_job]
import (
	"context"
	"fmt"
	"io"

	transcoder "cloud.google.com/go/video/transcoder/apiv1"
	"cloud.google.com/go/video/transcoder/apiv1/transcoderpb"
)

// getJob gets a previously-created job. See https://cloud.google.com/transcoder/docs/how-to/jobs#check_job_status
// for more information.
func getJob(w io.Writer, projectID string, location string, jobID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// jobID := "my-job-id"
	ctx := context.Background()
	client, err := transcoder.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &transcoderpb.GetJobRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectID, location, jobID),
	}

	response, err := client.GetJob(ctx, req)
	if err != nil {
		return fmt.Errorf("GetJob: %w", err)
	}

	fmt.Fprintf(w, "Job: %v", response)
	return nil
}

// [END transcoder_get_job]
