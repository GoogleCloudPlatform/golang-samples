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

// [START job_search_histogram_search]
import (
	"context"
	"fmt"
	"io"

	talent "cloud.google.com/go/talent/apiv4beta1"
	"cloud.google.com/go/talent/apiv4beta1/talentpb"
)

// histogramSearch searches for jobs with histogram queries.
func histogramSearch(w io.Writer, projectID, companyID string) error {
	ctx := context.Background()

	// Initialize a jobService client.
	c, err := talent.NewJobClient(ctx)
	if err != nil {
		return fmt.Errorf("talent.NewJobClient: %w", err)
	}
	defer c.Close()

	// Construct a searchJobs request.
	req := &talentpb.SearchJobsRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		// Make sure to set the RequestMetadata the same as the associated
		// search request.
		RequestMetadata: &talentpb.RequestMetadata{
			// Make sure to hash your userId.
			UserId: "HashedUsrID",
			// Make sure to hash the sessionId.
			SessionId: "HashedSessionID",
			// Domain of the website where the search is conducted.
			Domain: "www.googlesample.com",
		},
		HistogramQueries: []*talentpb.HistogramQuery{
			{
				// More info on histogram facets, constants, and built-in functions:
				// https://godoc.org/google.golang.org/genproto/googleapis/cloud/talent/v4beta1#SearchJobsRequest
				HistogramQuery: "count(base_compensation, [bucket(12, 20)])",
			},
		},
	}
	if companyID != "" {
		req.JobQuery = &talentpb.JobQuery{
			Companies: []string{fmt.Sprintf("projects/%s/companies/%s", projectID, companyID)},
		}
	}

	resp, err := c.SearchJobs(ctx, req)
	if err != nil {
		return fmt.Errorf("SearchJobs: %w", err)
	}

	for _, job := range resp.GetMatchingJobs() {
		fmt.Fprintf(w, "Job: %q\n", job.GetJob().GetName())
	}

	return nil
}

// [END job_search_histogram_search]
