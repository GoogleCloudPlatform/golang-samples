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

package howto

// [START job_search_list_jobs]
import (
	"context"
	"fmt"
	"io"

	talent "cloud.google.com/go/talent/apiv4beta1"
	"cloud.google.com/go/talent/apiv4beta1/talentpb"
	"google.golang.org/api/iterator"
)

// listJobs lists jobs with a filter, for example
// `companyName="projects/my-project/companies/123"`.
func listJobs(w io.Writer, projectID, companyID string) error {
	ctx := context.Background()

	// Initialize a jobService client.
	c, err := talent.NewJobClient(ctx)
	if err != nil {
		return fmt.Errorf("talent.NewJobClient: %w", err)
	}
	defer c.Close()

	// Construct a listJobs request.
	companyName := fmt.Sprintf("projects/%s/companies/%s", projectID, companyID)
	req := &talentpb.ListJobsRequest{
		Parent: "projects/" + projectID,
		Filter: fmt.Sprintf("companyName=%q", companyName),
	}

	it := c.ListJobs(ctx, req)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return fmt.Errorf("it.Next: %w", err)
		}
		fmt.Fprintf(w, "Listing job: %v\n", resp.GetName())
		fmt.Fprintf(w, "Job title: %v\n", resp.GetTitle())
	}
}

// [END job_search_list_jobs]
