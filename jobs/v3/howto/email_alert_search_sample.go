// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

import (
	"fmt"
	"io"
	"os"
	"time"

	talent "google.golang.org/api/jobs/v3"
)

// [START search_for_alerts]

// SearchForAlerts searches for jobs with email alert set which could receive updates later if search result updates
func searchForAlerts(service *talent.Service, parent string, companyName string) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata{
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
	}

	searchJobsRequest := &talent.SearchJobsRequest{
		RequestMetadata: requestMetadata,
		// Set the search mode to a regular search
		SearchMode: "JOB_SEARCH",
	}
	if companyName != "" {
		jobQuery := &talent.JobQuery{
			CompanyNames: []string{companyName},
		}
		searchJobsRequest.JobQuery = jobQuery
	}

	resp, err := service.Projects.Jobs.SearchForAlert(parent, searchJobsRequest).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with alerts: %v", err)
	}
	return resp, nil
}

// [END search_for_alerts]

// [START run_email_alert_search_sample]

func runEmailAlertSearchSample(w io.Writer) {
	parent := fmt.Sprintf("projects/%s", os.Getenv("GOOGLE_CLOUD_PROJECT"))
	service, _ := createCTSService()

	// Create a company before creating jobs
	companyToCreate := constructCompanyWithRequiredFields()
	companyCreated, _ := createCompany(service, parent, companyToCreate)
	fmt.Fprintf(w, "CreateCompany: %s\n", companyCreated.DisplayName)

	// Create a SDE job
	jobTitle := "Software Engineer"
	jobToCreate := constructJobWithRequiredFields(companyCreated.Name, jobTitle)
	jobCreated, _ := createJob(service, parent, jobToCreate)
	fmt.Fprintf(w, "CreateJob: %s\n", jobCreated.Title)

	// Wait for 10 seconds for post processing
	time.Sleep(10 * time.Second)

	// Search jobs with alerts
	resp, _ := searchForAlerts(service, parent, companyCreated.Name)
	fmt.Fprintf(w, "SearchForAlerts StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Fprintf(w, "-- match job: %s\n", mJob.Job.Title)
	}

	empty, _ := deleteJob(service, jobCreated.Name)
	fmt.Fprintf(w, "DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	empty, _ = deleteCompany(service, companyCreated.Name)
	fmt.Fprintf(w, "DeleteCompany StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)

}

// [END run_email_alert_search_sample]
