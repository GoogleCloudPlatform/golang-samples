// Copyright 2019 Google LLC
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

package jobs

// [START dlp_list_jobs]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"google.golang.org/api/iterator"
)

// listJobs lists jobs matching the given optional filter and optional jobType.
func listJobs(w io.Writer, projectID, filter, jobType string) error {
	// projectID := "my-project-id"
	// filter := "`state` = FINISHED"
	// jobType := "RISK_ANALYSIS_JOB"
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %w", err)
	}
	defer client.Close()

	// Create a configured request.
	req := &dlppb.ListDlpJobsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Filter: filter,
		Type:   dlppb.DlpJobType(dlppb.DlpJobType_value[jobType]),
	}
	// Send the request and iterate over the results.
	it := client.ListDlpJobs(ctx, req)
	for {
		j, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Next: %w", err)
		}
		fmt.Fprintf(w, "Job %v status: %v\n", j.GetName(), j.GetState())
	}
	return nil
}

// [END dlp_list_jobs]
