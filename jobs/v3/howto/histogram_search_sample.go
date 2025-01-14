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

import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	talent "google.golang.org/api/jobs/v3"
)

// [START job_histogram_search]

// histogramSearch searches for jobs with histogram facets.
func histogramSearch(w io.Writer, projectID, companyName string) (*talent.SearchJobsResponse, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %w", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %w", err)
	}

	parent := "projects/" + projectID
	req := &talent.SearchJobsRequest{
		// Make sure to set the RequestMetadata the same as the associated
		// search request.
		RequestMetadata: &talent.RequestMetadata{
			// Make sure to hash your userID.
			UserId: "HashedUsrId",
			// Make sure to hash the sessionID.
			SessionId: "HashedSessionId",
			// Domain of the website where the search is conducted.
			Domain: "www.googlesample.com",
		},
		HistogramFacets: &talent.HistogramFacets{
			SimpleHistogramFacets: []string{"COMPANY_ID"},
			CustomAttributeHistogramFacets: []*talent.CustomAttributeHistogramRequest{
				{
					Key:                  "someFieldString",
					StringValueHistogram: true,
				},
			},
		},
		// Set the search mode to a regular search.
		SearchMode:               "JOB_SEARCH",
		RequirePreciseResultSize: true,
	}
	if companyName != "" {
		req.JobQuery = &talent.JobQuery{
			CompanyNames: []string{companyName},
		}
	}

	resp, err := service.Projects.Jobs.Search(parent, req).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with Historgram Facets: %w", err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END job_histogram_search]
