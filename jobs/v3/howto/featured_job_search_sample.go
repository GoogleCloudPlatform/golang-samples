// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	talent "google.golang.org/api/jobs/v3"
)

// [START featured_job]

// constructFeaturedJob constructs a job as featured/promoted one
func constructFeaturedJob(companyName string, jobTitle string) *talent.Job {
	requisitionId := fmt.Sprintf("featured-job-required-fields-%d", time.Now().UnixNano())
	applicationInfo := &talent.ApplicationInfo{
		Uris: []string{"https://googlesample.com/career"},
	}
	job := &talent.Job{
		RequisitionId:   requisitionId,
		Title:           jobTitle,
		CompanyName:     companyName,
		ApplicationInfo: applicationInfo,
		Description:     "Design, devolop, test, deploy, maintain and improve software.",
		PromotionValue:  2,
	}
	//	fmt.Fprintf(w, "Job constructed: %v\n",job)
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
		log.Fatalf("Failed to search for jobs with query %v, Err: %v", query, err)
	}
	return resp, err
}

// [END search_featured_job]

// [START run_featured_job_search_sample]

func runFeaturedJobSearchSample(w io.Writer) {
	parent := fmt.Sprintf("projects/%s", os.Getenv("GOOGLE_CLOUD_PROJECT"))
	service, _ := createCtsService()

	// Create a company before creating jobs
	companyToCreate := constructCompanyWithRequiredFields()
	companyCreated, _ := createCompany(service, parent, companyToCreate)
	fmt.Fprintf(w, "CreateCompany: %s\n", companyCreated.DisplayName)

	// Construct a job
	jobTitle := "Software Engineer (Featured)"
	jobToCreate := constructFeaturedJob(companyCreated.Name, jobTitle)

	// Create a featured job
	jobCreated, _ := createJob(service, parent, jobToCreate)
	fmt.Fprintf(w, "CreateJob: %s\n", jobCreated.Title)

	// Wait for 10 seconds for post processing
	time.Sleep(10 * time.Second)

	// Search for featured jobs
	resp, _ := searchFeaturedJobs(service, parent, companyCreated.Name, jobTitle)
	fmt.Fprintf(w, "SearchFeaturedJobs StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Fprintf(w, "-- match job: %s\n", mJob.Job.Title)
	}

	// Delete Job
	empty, _ := deleteJob(service, jobCreated.Name)
	fmt.Fprintf(w, "DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	// Delete Company
	emptyResp, _ := deleteCompany(service, companyCreated.Name)
	fmt.Fprintf(w, "DeleteCompany StatusCode: %d\n", emptyResp.ServerResponse.HTTPStatusCode)
}

// [END run_featured_job_search_sample]
