// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package howto

import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	talent "google.golang.org/api/jobs/v3"
)

// [START basic_keyword_search]

// basicJobSearch searches for jobs with query.
func basicJobSearch(w io.Writer, projectID, companyName, query string) (*talent.SearchJobsResponse, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %v", err)
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
		return nil, fmt.Errorf("failed to search for jobs with query %q: %v", query, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END basic_keyword_search]

// [START category_filter]

// categoryFilterSearch searches for jobs on category filter.
func categoryFilterSearch(w io.Writer, projectID, companyName string, categories []string) (*talent.SearchJobsResponse, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %v", err)
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
		return nil, fmt.Errorf("failed to search for jobs with categories %v: %v", categories, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END category_filter]

// [START employment_types_filter]

// employmentTypesSearch searches for jobs on employment types.
func employmentTypesSearch(w io.Writer, projectID, companyName string, employmentTypes []string) (*talent.SearchJobsResponse, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %v", err)
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
		return nil, fmt.Errorf("failed to search for jobs with employment types %v: %v", employmentTypes, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END employment_types_filter]

// [START date_range_filter]

///dateRangeSearch searches for jobs on date range.
// In JSON format, the Timestamp type is encoded as a string in the
// [RFC 3339](https://www.ietf.org/rfc/rfc3339.txt) format. That is, the
// format is "{year}-{month}-{day}T{hour}:{min}:{sec}[.{frac_sec}]Z"
// e.g. "2017-01-15T01:30:15.01Z".
func dateRangeSearch(w io.Writer, projectID, companyName, startTime, endTime string) (*talent.SearchJobsResponse, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %v", err)
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
		return nil, fmt.Errorf("failed to search for jobs with date range [%s, %s]: %v", startTime, endTime, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END date_range_filter]

// [START language_code_filter]

// languageCodeSearch searches for jobs on language code.
func languageCodeSearch(w io.Writer, projectID, companyName string, languageCodes []string) (*talent.SearchJobsResponse, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %v", err)
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
		return nil, fmt.Errorf("failed to search for jobs with languange codes %v: %v", languageCodes, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END language_code_filter]

// [START company_display_name_filter]

// companyDisplayNameSearch searches for job on company display names.
func companyDisplayNameSearch(w io.Writer, projectID, companyName string, companyDisplayNames []string) (*talent.SearchJobsResponse, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %v", err)
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
		return nil, fmt.Errorf("failed to search for jobs with company display names %v: %v", companyDisplayNames, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END company_display_name_filter]

// [START compensation_fiter]

// compensationSearch searches for job on compensation.
func compensationSearch(w io.Writer, projectID, companyName string) (*talent.SearchJobsResponse, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %v", err)
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
		return nil, fmt.Errorf("failed to search for jobs with compensation: %v", err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END compensation_filter]
