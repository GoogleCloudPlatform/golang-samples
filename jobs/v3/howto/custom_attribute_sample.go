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

// [START custom_attribute_job]

// constructJobWithCustomAttributes constructs a job with custom attributes.
func constructJobWithCustomAttributes(companyName string, jobTitle string) *talent.Job {
	// requisitionID shoud be the unique ID in your system
	requisitionID := fmt.Sprintf("job-with-custom-attribute-%d", time.Now().UnixNano())
	applicationInfo := &talent.ApplicationInfo{
		Uris: []string{"https://googlesample.com/career"},
	}
	customAttrStr := &talent.CustomAttribute{
		Filterable:   true,
		StringValues: []string{"someStrVal"},
	}
	customAttrLong := &talent.CustomAttribute{
		Filterable: true,
		LongValues: []int64{900},
	}

	customAttributes := map[string]talent.CustomAttribute{
		"someFieldString": *customAttrStr,
		"someFieldLong":   *customAttrLong,
	}

	job := &talent.Job{
		RequisitionId:    requisitionID,
		Title:            jobTitle,
		CompanyName:      companyName,
		ApplicationInfo:  applicationInfo,
		Description:      "Design, devolop, test, deploy, maintain and improve software.",
		CustomAttributes: customAttributes,
	}
	return job
}

// [END custom_attribute_job]

// [START custom_attribute_filter_string_value]

// filterOnStringValueCustomAttribute searches for jobs on a string value custom atrribute.
func filterOnStringValueCustomAttribute(service *talent.Service, parent string) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata{
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
	}

	customAttrFilter := "NOT EMPTY(someFieldString)"
	query := &talent.JobQuery{
		CustomAttributeFilter: customAttrFilter,
	}
	searchJobsRequest := &talent.SearchJobsRequest{
		JobQuery:        query,
		RequestMetadata: requestMetadata,
		JobView:         "JOB_VIEW_FULL",
	}
	resp, err := service.Projects.Jobs.Search(parent, searchJobsRequest).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with string value custom attribute: %v", err)
	}

	return resp, nil
}

// [END custom_attribute_filter_string_value]

// [START custom_attribute_filter_long_value]

// filterOnLongValueCustomAttribute searches for jobs on a long value custom atrribute.
func filterOnLongValueCustomAttribute(service *talent.Service, parent string) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata{
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
	}

	customAttrFilter := "someFieldLong < 1000"
	query := &talent.JobQuery{
		CustomAttributeFilter: customAttrFilter,
	}
	searchJobsRequest := &talent.SearchJobsRequest{
		JobQuery:        query,
		RequestMetadata: requestMetadata,
		JobView:         "JOB_VIEW_FULL",
	}
	resp, err := service.Projects.Jobs.Search(parent, searchJobsRequest).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with long value custom attribute: %v", err)
	}

	return resp, nil
}

// [END custom_attribute_filter_long_value]

// [START custom_attribute_filter_multi_attributes]

// filterOnLongValueCustomAttribute searches for jobs on multiple custom atrributes.
func filterOnMultiCustomAttributes(service *talent.Service, parent string) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata{
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
	}

	customAttrFilter := "(someFieldString = \"someStrVal\") AND (someFieldLong < 1000)"
	query := &talent.JobQuery{
		CustomAttributeFilter: customAttrFilter,
	}

	searchJobsRequest := &talent.SearchJobsRequest{
		JobQuery:        query,
		RequestMetadata: requestMetadata,
		JobView:         "JOB_VIEW_FULL",
	}
	resp, err := service.Projects.Jobs.Search(parent, searchJobsRequest).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with multiple custom attributes: %v", err)
	}

	return resp, nil
}

// [END custom_attribute_filter_multi_attributes]

// [START run_custom_attribute_sample]

func runCustomAttributeSample(w io.Writer, projectID string) {
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

	// Create a job with custom fields
	jobTitle := "Software Engineer"
	jobToCreate := constructJobWithCustomAttributes(companyCreated.Name, jobTitle)
	jobCreated, err := createJob(service, parent, jobToCreate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "CreateJob: %s\n", jobCreated.Title)

	time.Sleep(10 * time.Second)

	resp, err := filterOnStringValueCustomAttribute(service, parent)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "FilterOnStringValueCustomAttribute StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Fprintf(w, "-- match job: %s\n", mJob.Job.Title)
	}

	resp, _ = filterOnLongValueCustomAttribute(service, parent)
	fmt.Fprintf(w, "FilterOnLongValueCustomAttribute StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Fprintf(w, "-- match job: %s\n", mJob.Job.Title)
	}

	resp, err = filterOnMultiCustomAttributes(service, parent)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "FilterOnMultiCustomAttributes StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
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
	empty, err = deleteCompany(service, companyCreated.Name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "DeleteCompany StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
}

// [END run_custom_attribute_sample]
