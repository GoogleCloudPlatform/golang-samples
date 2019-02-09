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

// [START basic_location_search]

// basicLocationSearch searches for jobs within distance of location.
func basicLocationSearch(w io.Writer, projectID, companyName, location string, distance float64) (*talent.SearchJobsResponse, error) {
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
		LocationFilters: []*talent.LocationFilter{
			{
				Address:         location,
				DistanceInMiles: distance,
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
		return nil, fmt.Errorf("failed to search for jobs with basic location %s within %f miles: %v", location, distance, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END basic_location_search]

// [START city_location_search]

// cityLocationSearch searches for jobs in the same city of given location.
func cityLocationSearch(w io.Writer, projectID, companyName, location string) (*talent.SearchJobsResponse, error) {
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
		LocationFilters: []*talent.LocationFilter{
			{
				Address: location,
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
		return nil, fmt.Errorf("failed to search for jobs with city location %s: %v", location, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END city_location_search]

// [START broadening_location_search]

// broadeningLocationSearch searches for jobs with a broadening area of given
// location.
func broadeningLocationSearch(w io.Writer, projectID, companyName, location string) (*talent.SearchJobsResponse, error) {
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
		LocationFilters: []*talent.LocationFilter{
			{
				Address: location,
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
		SearchMode:       "JOB_SEARCH",
		EnableBroadening: true,
	}
	resp, err := service.Projects.Jobs.Search(parent, req).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with broadening location %v: %v", location, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END broadening_location_search]

// [START keyword_location_search]

// keywordLocationSearch searches for jobs with given keyword and within the
// distance of given location.
func keywordLocationSearch(w io.Writer, projectID, companyName, location string, distance float64, keyword string) (*talent.SearchJobsResponse, error) {
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
		LocationFilters: []*talent.LocationFilter{
			{
				Address:         location,
				DistanceInMiles: distance,
			},
		},
		Query: keyword,
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
		return nil, fmt.Errorf("failed to search for jobs with keyword %q in location %v within %f miles: %v", keyword, location, distance, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END keyword_location_search]

// [START multi_locations_search]

// multiLocationsSearch searches for jobs that fall in the distance of any given
// locations.
func multiLocationsSearch(w io.Writer, projectID, companyName, location, location2 string, distance float64) (*talent.SearchJobsResponse, error) {
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
		LocationFilters: []*talent.LocationFilter{
			{
				Address:         location,
				DistanceInMiles: distance,
			},
			{
				Address:         location2,
				DistanceInMiles: distance,
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
		return nil, fmt.Errorf("Failed to search for jobs with multi locations %s and %s within %f miles, Err: %v", location, location2, distance, err)
	}

	fmt.Fprintln(w, "Jobs:")
	for _, j := range resp.MatchingJobs {
		fmt.Fprintf(w, "\t%q\n", j.Job.Name)
	}

	return resp, nil
}

// [END multi_locations_search]
