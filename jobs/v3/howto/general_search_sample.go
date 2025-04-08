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

// [START job_discovery_basic_keyword_search]

// basicJobSearch searches for jobs with query.
func basicJobSearch(w io.Writer, projectID, companyName, query string) (*talent.SearchJobsResponse, error) {
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
		// Make sure to set the requestMetadata the same as the associated
		// search request.
		RequestMetadata: &talent.RequestMetadata{
			// Make sure to hash your userID
			UserId: "HashedUsrId",
			// Make sure to hash the sessionID
			SessionId: "HashedSessionId",
			// Domain of the website where the search is conducted
			Domain: "www.googlesample.com",
		},
		// Set the actual search term as defined in the jobQuery
		JobQuery: jobQuery,
		// Set the search mode to a regular search
		SearchMode: "JOB_SEARCH",
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

// [END job_discovery_basic_keyword_search]

// [START job_discovery_category_filter_search]

// categoryFilterSearch searches for jobs on category filter.
func categoryFilterSearch(w io.Writer, projectID, companyName string, categories []string) (*talent.SearchJobsResponse, error) {
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
		JobCategories: categories,
	}
	if companyName != "" {
		jobQuery.CompanyNames = []string{companyName}
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
		// Set the actual search term as defined in the jobQuery.
		JobQuery: jobQuery,
		// Set the search mode to a regular search.
		SearchMode: "JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(parent, req).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with categories %v: %w", categories, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END job_discovery_category_filter_search]

// [START job_discovery_employment_types_filter_search]

// employmentTypesSearch searches for jobs on employment types.
func employmentTypesSearch(w io.Writer, projectID, companyName string, employmentTypes []string) (*talent.SearchJobsResponse, error) {
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
		EmploymentTypes: employmentTypes,
	}
	if companyName != "" {
		jobQuery.CompanyNames = []string{companyName}
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
		// Set the actual search term as defined in the jobQuery.
		JobQuery: jobQuery,
		// Set the search mode to a regular search.
		SearchMode: "JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(parent, req).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with employment types %v: %w", employmentTypes, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END job_discovery_employment_types_filter_search]

// [START job_discovery_date_range_filter_search]

// /dateRangeSearch searches for jobs on date range.
// In JSON format, the Timestamp type is encoded as a string in the
// [RFC 3339](https://www.ietf.org/rfc/rfc3339.txt) format. That is, the
// format is "{year}-{month}-{day}T{hour}:{min}:{sec}[.{frac_sec}]Z"
// e.g. "2017-01-15T01:30:15.01Z".
func dateRangeSearch(w io.Writer, projectID, companyName, startTime, endTime string) (*talent.SearchJobsResponse, error) {
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
		PublishTimeRange: &talent.TimestampRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
	}
	if companyName != "" {
		jobQuery.CompanyNames = []string{companyName}
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
		// Set the actual search term as defined in the jobQuery.
		JobQuery: jobQuery,
		// Set the search mode to a regular search.
		SearchMode: "JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(parent, req).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with date range [%s, %s]: %w", startTime, endTime, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END job_discovery_date_range_filter_search]

// [START job_discovery_language_code_filter_search]

// languageCodeSearch searches for jobs on language code.
func languageCodeSearch(w io.Writer, projectID, companyName string, languageCodes []string) (*talent.SearchJobsResponse, error) {
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
		LanguageCodes: languageCodes,
	}
	if companyName != "" {
		jobQuery.CompanyNames = []string{companyName}
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
		// Set the actual search term as defined in the jobQuery.
		JobQuery: jobQuery,
		// Set the search mode to a regular search.
		SearchMode: "JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(parent, req).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with languange codes %v: %w", languageCodes, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END job_discovery_language_code_filter_search]

// [START job_discovery_company_display_name_search]

// companyDisplayNameSearch searches for job on company display names.
func companyDisplayNameSearch(w io.Writer, projectID, companyName string, companyDisplayNames []string) (*talent.SearchJobsResponse, error) {
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
		CompanyDisplayNames: companyDisplayNames,
	}
	if companyName != "" {
		jobQuery.CompanyNames = []string{companyName}
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
		// Set the actual search term as defined in the jobQuery.
		JobQuery: jobQuery,
		// Set the search mode to a regular search.
		SearchMode: "JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(parent, req).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with company display names %v: %w", companyDisplayNames, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END job_discovery_company_display_name_search]

// [START job_discovery_compensation_search]

// compensationSearch searches for job on compensation.
func compensationSearch(w io.Writer, projectID, companyName string) (*talent.SearchJobsResponse, error) {
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
		CompensationFilter: &talent.CompensationFilter{
			Type:  "UNIT_AND_AMOUNT",
			Units: []string{"HOURLY"},
			Range: &talent.CompensationRange{
				MaxCompensation: &talent.Money{
					Units:        15,
					CurrencyCode: "USD",
				},
				MinCompensation: &talent.Money{
					Units:        10,
					CurrencyCode: "USD",
					Nanos:        500000000,
				},
			},
		},
	}
	if companyName != "" {
		jobQuery.CompanyNames = []string{companyName}
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
		// Set the actual search term as defined in the jobQuery.
		JobQuery: jobQuery,
		// Set the search mode to a regular search.
		SearchMode: "JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(parent, req).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with compensation: %w", err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END job_discovery_compensation_search]
