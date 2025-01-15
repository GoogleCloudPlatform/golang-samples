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
	"time"

	"golang.org/x/oauth2/google"
	talent "google.golang.org/api/jobs/v3"
)

// [START job_featured_job]

// constructFeaturedJob constructs a job as featured/promoted one.
func constructFeaturedJob(companyName string, jobTitle string) *talent.Job {
	requisitionID := fmt.Sprintf("featured-job-required-fields-%d", time.Now().UnixNano())
	applicationInfo := &talent.ApplicationInfo{
		Uris: []string{"https://googlesample.com/career"},
	}
	job := &talent.Job{
		RequisitionId:   requisitionID,
		Title:           jobTitle,
		CompanyName:     companyName,
		ApplicationInfo: applicationInfo,
		Description:     "Design, devolop, test, deploy, maintain and improve software.",
		PromotionValue:  2,
	}
	return job
}

// [END job_featured_job]

// [START job_search_featured_job]

// searchFeaturedJobs searches for jobs with query.
func searchFeaturedJobs(w io.Writer, projectID, companyName, query string) (*talent.SearchJobsResponse, error) {
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

	jobQuery := &talent.JobQuery{
		Query: query,
	}
	if companyName != "" {
		jobQuery.CompanyNames = []string{companyName}
	}

	parent := "projects/" + projectID
	req := &talent.SearchJobsRequest{
		// Make sure to set the RequestMetadata the same as the associated
		// Search request.
		RequestMetadata: &talent.RequestMetadata{
			// Make sure to hash your userID.
			UserId: "HashedUsrId",
			// Make sure to hash the sessionID.
			SessionId: "HashedSessionId",
			// Domain of the website where the search is conducted.
			Domain: "www.googlesample.com",
		},
		// Set the actual search term as defined in the jobQuery.
		JobQuery: jobQuery,
		// Set the search mode to a featured search, wwhich only searches for
		// jobs with a positive promotion value.
		SearchMode: "FEATURED_JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(parent, req).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with query %q: %w", query, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END job_search_featured_job]
