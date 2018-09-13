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

// [START commute_search]

// commuteSearch searches for jobs within commute filter.
func commuteSearch(service *talent.Service, parent string, companyName string) (*talent.SearchJobsResponse, error) {
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
		CommuteFilter: &talent.CommuteFilter{
			RoadTraffic:    "TRAFFIC_FREE",
			CommuteMethod:  "TRANSIT",
			TravelDuration: "1000s",
			StartCoordinates: &talent.LatLng{
				Latitude:  37.422408,
				Longitude: -122.085609,
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
		SearchMode:               "JOB_SEARCH",
		RequirePreciseResultSize: true,
	}
	resp, err := service.Projects.Jobs.Search(parent, searchJobsRequest).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with commute filter: %v", err)
	}
	return resp, nil
}

// [END commute_search]

// [START run_commute_search_sample]
func runCommuteSearchSample(w io.Writer, projectID string) {
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

	jobTitle := "Software Engineer"
	jobToCreate := constructJobWithRequiredFields(companyCreated.Name, jobTitle)
	jobToCreate.Addresses = []string{"1600 Amphitheatre Parkway, Mountain View, CA 94043"}
	jobCreated, err := createJob(service, parent, jobToCreate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "CreateJob: %s\n", jobCreated.Title)

	// Wait several seconds for post processing
	time.Sleep(10 * time.Second)

	resp, err := commuteSearch(service, parent, companyCreated.Name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "CommuteSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Fprintf(w, "-- match job: %s\n", mJob.Job.Title)
	}

	empty, err := deleteJob(service, jobCreated.Name)
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

// [END run_commute_search_sample]
