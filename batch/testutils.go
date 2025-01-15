// Copyright 2022 Google LLC
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

package snippets

import (
	"context"
	"fmt"
	"time"

	batch "cloud.google.com/go/batch/apiv1"
	"cloud.google.com/go/batch/apiv1/batchpb"
)

func jobSucceeded(projectID, region, jobName string) (bool, error) {
	ctx := context.Background()
	batchClient, err := batch.NewClient(ctx)
	if err != nil {
		return false, fmt.Errorf("NewClient: %w", err)
	}
	defer batchClient.Close()

	const maxAttempts = 30

	for i := 0; i < 120; i++ {
		req := &batchpb.GetJobRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectID, region, jobName),
		}
		response, err := batchClient.GetJob(ctx, req)
		if err != nil {
			return false, fmt.Errorf("unable to get job: %w", err)
		}
		if response.GetStatus().State == batchpb.JobStatus_SUCCEEDED {
			return true, nil
		}
		if response.GetStatus().State == batchpb.JobStatus_FAILED {
			return false, nil
		}
		time.Sleep(2 * time.Second)
	}
	return false, fmt.Errorf("Timed out waiting for job to succeed or fail after %d attempts", maxAttempts)
}
