// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

import (
	"fmt"
	"io"
	"log"
	"time"

	talent "google.golang.org/api/jobs/v3"
)

// [START basic_location_search]

// basicLocationSearch searches for jobs within distance of location
func basicLocationSearch(service *talent.Service, parent string, companyName string, location string, distance float64) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata{
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
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

	searchJobsRequest := &talent.SearchJobsRequest{
		RequestMetadata: requestMetadata,
		// Set the actual search term as defined in the jobQurey
		JobQuery: jobQuery,
		// Set the search mode to a regular search
		SearchMode: "JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(parent, searchJobsRequest).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with basic location %s within %f miles: %v", location, distance, err)
	}
	return resp, nil
}

// [END basic_location_search]

// [START city_location_search]

// cityLocationSearch searches for jobs in the same city of given location.
func cityLocationSearch(service *talent.Service, parent string, companyName string, location string) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata{
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
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

	searchJobsRequest := &talent.SearchJobsRequest{
		RequestMetadata: requestMetadata,
		// Set the actual search term as defined in the jobQurey
		JobQuery: jobQuery,
		// Set the search mode to a regular search
		SearchMode: "JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(parent, searchJobsRequest).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with city location %s: %v", location, err)
	}
	return resp, nil
}

// [END city_location_search]

// [START broadening_location_search]

// broadeningLocationSearch searches for jobs with a broadening area of given location.
func broadeningLocationSearch(service *talent.Service, parent string, companyName string, location string) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata{
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
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

	searchJobsRequest := &talent.SearchJobsRequest{
		RequestMetadata: requestMetadata,
		// Set the actual search term as defined in the jobQurey
		JobQuery: jobQuery,
		// Set the search mode to a regular search
		SearchMode:       "JOB_SEARCH",
		EnableBroadening: true,
	}
	resp, err := service.Projects.Jobs.Search(parent, searchJobsRequest).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with broadening location %v: %v", location, err)
	}
	return resp, nil
}

// [END broadening_location_search]

// [START keyword_location_search]

// keywordLocationSearch searches for jobs with given keyword and within the distance of given location.
func keywordLocationSearch(service *talent.Service, parent string, companyName string, location string, distance float64, keyword string) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata{
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
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

	searchJobsRequest := &talent.SearchJobsRequest{
		RequestMetadata: requestMetadata,
		// Set the actual search term as defined in the jobQurey
		JobQuery: jobQuery,
		// Set the search mode to a regular search
		SearchMode: "JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(parent, searchJobsRequest).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with keyword %q in location %v within %f miles: %v", keyword, location, distance, err)
	}
	return resp, nil
}

// [END keyword_location_search]

// [START multi_locations_search]

// multiLocationsSearch searches for jobs that fall in the distance of any given locations.
func multiLocationsSearch(service *talent.Service, parent string, companyName string, location string, distance float64, location2 string) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata{
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
	}
	jobQuery := &talent.JobQuery{
		LocationFilters: []*talent.LocationFilter{
			{
				Address:         location,
				DistanceInMiles: distance,
			},
			{
				Address: location2,
			},
		},
	}

	if companyName != "" {
		jobQuery.CompanyNames = []string{companyName}
	}

	searchJobsRequest := &talent.SearchJobsRequest{
		RequestMetadata: requestMetadata,
		// Set the actual search term as defined in the jobQurey
		JobQuery: jobQuery,
		// Set the search mode to a regular search
		SearchMode: "JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(parent, searchJobsRequest).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to search for jobs with multi locations %s and %s within %f miles, Err: %v", location, location2, distance, err)
	}
	return resp, nil
}

// [END multi_locations_search]

// [START run_location_based_search_sample]

func runLocationBasedSearchSample(w io.Writer, projectID string) {
	parent := fmt.Sprintf("projects/%s", projectID)
	service, err := createCTSService()
	if err != nil {
		log.Fatal(err)
	}

	// Create a company before creating jobs
	companyToCreate := constructCompanyWithRequiredFields()
	companyCreated, err := createCompany(service, parent, companyToCreate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "CreateCompany: %s\n", companyCreated.DisplayName)

	location := "Mountain View, CA"
	distance := 0.5
	keyword := "Software Engineer"
	location2 := "Sunnyvale, CA"

	// Create a SDE job
	jobTitle := keyword
	jobToCreate := constructJobWithRequiredFields(companyCreated.Name, jobTitle)
	jobToCreate.Addresses = []string{location}

	jobCreated, err := createJob(service, parent, jobToCreate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "CreateJob: %s\n", jobCreated.Title)

	jobTitle2 := "Senior " + keyword
	jobToCreate2 := constructJobWithRequiredFields(companyCreated.Name, jobTitle2)
	jobToCreate2.Addresses = []string{location2}
	jobCreated2, err := createJob(service, parent, jobToCreate2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "CreateJob: %s\n", jobCreated2.Title)

	// Wait for 10 seconds for post processing
	time.Sleep(10 * time.Second)

	resp, err := basicLocationSearch(service, parent, companyCreated.Name, location, distance)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "BasicLocationSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))

	resp, err = cityLocationSearch(service, parent, companyCreated.Name, location)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "CityLocationSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Fprintf(w, "-- match job: %s\n", mJob.Job.Title)
	}

	resp, err = broadeningLocationSearch(service, parent, companyCreated.Name, location)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "BroadeningLocationSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))

	resp, err = keywordLocationSearch(service, parent, companyCreated.Name, location, distance, keyword)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "KeywordLocationSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))

	resp, err = multiLocationsSearch(service, parent, companyCreated.Name, location, distance, location2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "MultiLocationsSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))

	empty, err := deleteJob(service, jobCreated.Name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	empty, err = deleteJob(service, jobCreated2.Name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	empty, err = deleteCompany(service, companyCreated.Name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "DeleteCompany StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
}

// [END run_location_based_search_sample]
