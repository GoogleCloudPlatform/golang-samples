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

// [START auto_complete_job_title]

// jobTitleAutoComplete suggests the job titles of the given companyName based on query.
func jobTitleAutoComplete(service *talent.Service, parent string, companyName string, query string) (*talent.CompleteQueryResponse, error) {
	complete := service.Projects.Complete(parent).Query(query).LanguageCode("en-US").Type("JOB_TITLE").PageSize(10)
	if companyName != "" {
		complete.CompanyName(companyName)
	}
	resp, err := complete.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to auto complete with query %s in company %s: %v", query, companyName, err)
	}
	return resp, nil

}

// [END auto_complete_job_title]

// [START auto_complete_default]

// defaultAutoComplete suggests job titles or company display names of given companyName based on query.
func defaultAutoComplete(service *talent.Service, parent string, companyName string, query string) (*talent.CompleteQueryResponse, error) {
	complete := service.Projects.Complete(parent).Query(query).LanguageCode("en-US").Type("COMBINED").PageSize(10)
	if companyName != "" {
		complete.CompanyName(companyName)
	}
	resp, err := complete.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to auto complete with query %s in company %s: %v", query, companyName, err)
	}
	return resp, nil

}

// [END auto_complete_default]

// [START run_auto_complete_sample]

// runAutoCompleteSample is to run all samples of auto complete APIs.
func runAutoCompleteSample(w io.Writer) {
	parent := fmt.Sprintf("projects/%s", os.Getenv("GOOGLE_CLOUD_PROJECT"))
	service, _ := createCTSService()

	// Create a company before creating jobs
	companyToCreate := constructCompanyWithRequiredFields()
	companyCreated, _ := createCompany(service, parent, companyToCreate)
	fmt.Fprintf(w, "CreateCompany: %s\n", companyCreated.DisplayName)

	// Create a SDE job
	jobTitleSWE := "Software Engineer"
	jobToCreateSWE := constructJobWithRequiredFields(companyCreated.Name, jobTitleSWE)
	jobCreatedSWE, _ := createJob(service, parent, jobToCreateSWE)
	fmt.Fprintf(w, "CreateJob: %s\n", jobCreatedSWE.Title)
	// Create a PM job
	jobTitlePM := "GAP Product Manager"
	jobToCreatePM := constructJobWithRequiredFields(companyCreated.Name, jobTitlePM)
	jobCreatedPM, _ := createJob(service, parent, jobToCreatePM)
	fmt.Fprintf(w, "CreateJob: %s\n", jobCreatedPM.Title)

	// Wait several seconds for post processing
	time.Sleep(10 * time.Second)

	query := "sof"
	resp, _ := defaultAutoComplete(service, parent, "" /*companyName*/, query)
	fmt.Fprintf(w, "DefaultAutoComplete query: %s StatusCode: %d\n", query, resp.ServerResponse.HTTPStatusCode)
	for _, comp := range resp.CompletionResults {
		fmt.Fprintf(w, "-- suggestion: %s\n", comp.Suggestion)
	}

	query = "sof"
	resp, _ = jobTitleAutoComplete(service, parent, "" /*companyName*/, query)
	fmt.Fprintf(w, "JobTitleAutoComplete query: %s StatusCode: %d\n", query, resp.ServerResponse.HTTPStatusCode)
	for _, comp := range resp.CompletionResults {
		fmt.Fprintf(w, "-- suggestion: %s\n", comp.Suggestion)
	}

	query = "gap"
	resp, _ = defaultAutoComplete(service, parent, companyCreated.Name, query)
	fmt.Fprintf(w, "DefaultAutoComplete query: %s StatusCode: %d\n", query, resp.ServerResponse.HTTPStatusCode)
	for _, comp := range resp.CompletionResults {
		fmt.Fprintf(w, "-- suggestion: %s\n", comp.Suggestion)
	}

	query = "gap"
	resp, _ = jobTitleAutoComplete(service, parent, companyCreated.Name, query)
	fmt.Fprintf(w, "JobTitleAutoComplete query: %s StatusCode: %d\n", query, resp.ServerResponse.HTTPStatusCode)
	for _, comp := range resp.CompletionResults {
		fmt.Fprintf(w, "-- suggestion: %s\n", comp.Suggestion)
	}

	// Delete Job
	empty, _ := deleteJob(service, jobCreatedSWE.Name)
	fmt.Fprintf(w, "DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	// Delete Job
	empty, _ = deleteJob(service, jobCreatedPM.Name)
	fmt.Fprintf(w, "DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	// Delete Company
	empty, _ = deleteCompany(service, companyCreated.Name)
	fmt.Fprintf(w, "DeleteCompany StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)

}

// [END run_auto_complete_sample]
