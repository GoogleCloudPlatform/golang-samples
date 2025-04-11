// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jobs

// [START dlp_get_job]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// jobsGet gets an inspection job using jobName
func jobsGet(w io.Writer, projectID string, jobName string) error {
	// projectId := "my-project-id"
	// jobName := "your-job-id"

	ctx := context.Background()

	// Initialize a client once and reuse it to send multiple requests. Clients
	// are safe to use across goroutines. When the client is no longer needed,
	// call the Close method to cleanup its resources.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}

	// Closing the client safely cleans up background resources.
	defer client.Close()

	// Construct the request to be sent by the client.
	req := &dlppb.GetDlpJobRequest{
		Name: jobName,
	}

	// Send the request.
	resp, err := client.GetDlpJob(ctx, req)
	if err != nil {
		return err
	}

	// Print the results.
	fmt.Fprintf(w, "Job Name: %v Job Status: %v", resp.Name, resp.State)
	return nil
}

// [END dlp_get_job]
