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

// [START featured_job]

// constructFeaturedJob constructs a job as featured/promoted one
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

// [END featured_job]

// [START search_featured_job]

// searchFeaturedJobs searches for jobs with query.
func searchFeaturedJobs(service *talent.Service, parent string, companyName string, query string) (*talent.SearchJobsResponse, error) {
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
		Query: query,
	}
	if companyName != "" {
		jobQuery.CompanyNames = []string{companyName}
	}

	searchJobsRequest := &talent.SearchJobsRequest{
		RequestMetadata: requestMetadata,
		// Set the actual search term as defined in the jobQurey
		JobQuery: jobQuery,
		// Set the search mode to a featured search
		// which would only search the jobs with positive promotion value
		SearchMode: "FEATURED_JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(parent, searchJobsRequest).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with query %q: %v", query, err)
	}
	return resp, nil
}

// [END search_featured_job]

// [START run_featured_job_search_sample]

func runFeaturedJobSearchSample(w io.Writer, projectID string) {
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

	// Construct a job
	jobTitle := "Software Engineer (Featured)"
	jobToCreate := constructFeaturedJob(companyCreated.Name, jobTitle)

	// Create a featured job
	jobCreated, err := createJob(service, parent, jobToCreate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "CreateJob: %s\n", jobCreated.Title)

	// Wait for 10 seconds for post processing
	time.Sleep(10 * time.Second)

	// Search for featured jobs
	resp, err := searchFeaturedJobs(service, parent, companyCreated.Name, jobTitle)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "SearchFeaturedJobs StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Fprintf(w, "-- match job: %s\n", mJob.Job.Title)
	}

	// Delete Job
	empty, err := deleteJob(service, jobCreated.Name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	// Delete Company
	emptyResp, err := deleteCompany(service, companyCreated.Name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "DeleteCompany StatusCode: %d\n", emptyResp.ServerResponse.HTTPStatusCode)
}

// [END run_featured_job_search_sample]
