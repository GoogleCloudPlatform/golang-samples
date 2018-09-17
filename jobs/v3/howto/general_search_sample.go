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

// [START basic_keyword_search]

// basicJobSearch searches for jobs with query.
func basicJobSearch(service *talent.Service, parent string, companyName string, query string) (*talent.SearchJobsResponse, error) {
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
		// Set the search mode to a regular search
		SearchMode: "JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(parent, searchJobsRequest).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with query %q: %v", query, err)
	}
	return resp, nil
}

// [END basic_keyword_search]

// [START category_filter]

// categoryFilterSearch searches for jobs on category filter.
func categoryFilterSearch(service *talent.Service, parent string, companyName string, categories []string) (*talent.SearchJobsResponse, error) {
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
		JobCategories: categories,
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
		return nil, fmt.Errorf("failed to search for jobs with categories %v: %v", categories, err)
	}
	return resp, nil
}

// [END category_filter]

// [START employment_types_filter]

// employmentTypesSearch searches for jobs on employment types.
func employmentTypesSearch(service *talent.Service, parent string, companyName string, employmentTypes []string) (*talent.SearchJobsResponse, error) {
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
		EmploymentTypes: employmentTypes,
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
		return nil, fmt.Errorf("failed to search for jobs with employment types %v: %v", employmentTypes, err)
	}
	return resp, nil
}

// [END employment_types_filter]

// [START date_range_filter]

/**
 * SdateRangeSearch searches for jobs on date range.
 * In JSON format, the Timestamp type is encoded as a string in the
 * [RFC 3339](https://www.ietf.org/rfc/rfc3339.txt) format. That is, the
 * format is "{year}-{month}-{day}T{hour}:{min}:{sec}[.{frac_sec}]Z"
 * e.g. "2017-01-15T01:30:15.01Z"
 */
func dateRangeSearch(service *talent.Service, parent string, companyName string, startTime string, endTime string) (*talent.SearchJobsResponse, error) {
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
		PublishTimeRange: &talent.TimestampRange{
			StartTime: startTime,
			EndTime:   endTime,
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
		return nil, fmt.Errorf("failed to search for jobs with date range [%s, %s]: %v", startTime, endTime, err)
	}
	return resp, nil
}

// [END date_range_filter]

// [START language_code_filter]

// languageCodeSearch searches for jobs on language code.
func languageCodeSearch(service *talent.Service, parent string, companyName string, languageCodes []string) (*talent.SearchJobsResponse, error) {
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
		LanguageCodes: languageCodes,
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
		return nil, fmt.Errorf("failed to search for jobs with languange codes %v: %v", languageCodes, err)
	}
	return resp, nil
}

// [END language_code_filter]

// [START company_display_name_filter]

// companyDisplayNameSearch searches for job on company display names
func companyDisplayNameSearch(service *talent.Service, parent string, companyName string, companyDisplayNames []string) (*talent.SearchJobsResponse, error) {
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
		CompanyDisplayNames: companyDisplayNames,
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
		return nil, fmt.Errorf("failed to search for jobs with company display names %v: %v", companyDisplayNames, err)
	}
	return resp, nil
}

// [END company_display_name_filter]

// [START compensation_fiter]

// compensationSearch searches for job on compensation
func compensationSearch(service *talent.Service, parent string, companyName string) (*talent.SearchJobsResponse, error) {
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

	searchJobsRequest := &talent.SearchJobsRequest{
		RequestMetadata: requestMetadata,
		// Set the actual search term as defined in the jobQurey
		JobQuery: jobQuery,
		// Set the search mode to a regular search
		SearchMode: "JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(parent, searchJobsRequest).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for jobs with compensation: %v", err)
	}
	return resp, nil
}

// [END compensation_filter]

// [START run_general_search_sample]

// runGeneralSearchSample runs the general job search samples
func runGeneralSearchSample(w io.Writer, projectID string) {
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

	// Create a job
	jobTitle := "Systems Administrator"
	jobToCreate := constructJobWithRequiredFields(companyCreated.Name, jobTitle)
	jobToCreate.LanguageCode = "en-US"
	jobToCreate.EmploymentTypes = []string{"FULL_TIME"}
	compensationInfo := &talent.CompensationInfo{}
	compensationInfo.Entries = []*talent.CompensationEntry{
		{
			Type: "BASE",
			Unit: "HOURLY",
			Amount: &talent.Money{
				CurrencyCode: "USD",
				Units:        12,
			},
		},
	}
	jobToCreate.CompensationInfo = compensationInfo

	jobCreated, err := createJob(service, parent, jobToCreate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "CreateJob: %s\n", jobCreated.Title)

	// Wait for 10 seconds for post processing
	time.Sleep(10 * time.Second)

	resp, err := basicJobSearch(service, parent, companyCreated.Name, jobTitle)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "BasicJobSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Fprintf(w, "-- match job: %s\n", mJob.Job.Title)
	}

	categories := []string{"COMPUTER_AND_IT"}
	resp, err = categoryFilterSearch(service, parent, companyCreated.Name, categories)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "CategoryFilterSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Fprintf(w, "-- match job: %s\n", mJob.Job.Title)
	}

	employmentTypes := []string{"FULL_TIME", "CONTRACTOR", "PER_DIEM"}
	resp, err = employmentTypesSearch(service, parent, companyCreated.Name, employmentTypes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "EmploymentTypesSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Fprintf(w, "-- match job: %s\n", mJob.Job.Title)
	}

	resp, err = dateRangeSearch(service, parent, companyCreated.Name, "2000-01-01T00:00:00.01Z", "2099-01-01T00:00:00.01Z")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "DateRangeSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Fprintf(w, "-- match job: %s\n", mJob.Job.Title)
	}

	resp, err = languageCodeSearch(service, parent, companyCreated.Name, []string{"pt-BR", "en-US"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "LanguageCodeSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Fprintf(w, "-- match job: %s\n", mJob.Job.Title)
	}

	resp, err = companyDisplayNameSearch(service, parent, companyCreated.Name, []string{"Google Sample"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "CompanyDisplayNameSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Fprintf(w, "MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Fprintf(w, "-- match job: %s\n", mJob.Job.Title)
	}

	resp, err = compensationSearch(service, parent, companyCreated.Name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "CompensationSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
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

// [END run_general_search_sample]
