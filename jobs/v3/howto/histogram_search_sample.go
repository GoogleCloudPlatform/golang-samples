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

// [START histogram_search]

// histogramSearch searches for jobs with histogram facets
func histogramSearch(service *talent.Service, parent string, companyName string) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata{
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
	}

	histogramFacets := &talent.HistogramFacets{
		SimpleHistogramFacets: []string{"COMPANY_ID"},
		CustomAttributeHistogramFacets: []*talent.CustomAttributeHistogramRequest{
			{
				Key:                  "someFieldString",
				StringValueHistogram: true,
			},
		},
	}

	searchJobsRequest := &talent.SearchJobsRequest{
		RequestMetadata: requestMetadata,
		HistogramFacets: histogramFacets,
		// Set the search mode to a regular search
		SearchMode:               "JOB_SEARCH",
		RequirePreciseResultSize: true,
	}
	if companyName != "" {
		jobQuery := &talent.JobQuery{
			CompanyNames: []string{companyName},
		}
		searchJobsRequest.JobQuery = jobQuery
	}

	resp, err := service.Projects.Jobs.Search(parent, searchJobsRequest).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with Historgram Facets: %v", err)
	}
	return resp, nil
}

// [END histogram_search]

// [START run_histogram_search_sample]

func runHistogramSearchSample(w io.Writer, projectID string) {
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

	// Create a SDE job
	jobTitleSWE := "Software Engineer"
	jobToCreateSWE := constructJobWithCustomAttributes(companyCreated.Name, jobTitleSWE)
	jobCreatedSWE, err := createJob(service, parent, jobToCreateSWE)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "CreateJob: %s\n", jobCreatedSWE.Title)

	// Wait several seconds for post processing
	time.Sleep(10 * time.Second)

	resp, err := histogramSearch(service, parent, companyCreated.Name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "HistogramSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Fprintf(w, "-- match job: %s\n", mJob.Job.Title)
	}
	fmt.Fprintf(w, "SimpleHistogramResults size: %d\n", len(resp.HistogramResults.SimpleHistogramResults))
	for _, hist := range resp.HistogramResults.SimpleHistogramResults {
		fmt.Fprintf(w, "-- simple histogram searchType: %s value: %v\n", hist.SearchType, hist.Values)
	}
	fmt.Fprintf(w, "CustomAttributeHistogramResults size: %d\n", len(resp.HistogramResults.CustomAttributeHistogramResults))
	for _, hist := range resp.HistogramResults.CustomAttributeHistogramResults {
		fmt.Fprintf(w, "-- custom-attribute histogram key: %s value: %v\n", hist.Key, hist.StringValueHistogramResult)
	}

	empty, err := deleteJob(service, jobCreatedSWE.Name)
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

// [END run_histogram_search_sample]
