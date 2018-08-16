package cjdsample

import (
	"fmt"
	"log"
	"time"

	talent "google.golang.org/api/jobs/v3"
)

// [START basic_keyword_search]

/**
 * Simple search jobs with keyword.
 */
func BasicJobSearch(service *talent.Service, companyName string, query string) (*talent.SearchJobsResponse, error) {
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
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with query %v, Err: %v", query, err)
	}
	return resp, err
}

// [END basic_keyword_search]

// [START category_filter]

/**
 * Search on category filter
 */
func CategoryFilterSearch(service *talent.Service, companyName string, categories []string) (*talent.SearchJobsResponse, error) {
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
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with categories %v, Err: %v", categories, err)
	}
	return resp, err
}

// [END category_filter]

// [START employment_types_filter]

/**
 * Search on employment types
 */
func EmploymentTypesSearch(service *talent.Service, companyName string, employmentTypes []string) (*talent.SearchJobsResponse, error) {
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
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with employment types %v, Err: %v", employmentTypes, err)
	}
	return resp, err
}

// [END employment_types_filter]

// [START date_range_filter]

/**
 * Search on date range.
 * In JSON format, the Timestamp type is encoded as a string in the
 * [RFC 3339](https://www.ietf.org/rfc/rfc3339.txt) format. That is, the
 * format is "{year}-{month}-{day}T{hour}:{min}:{sec}[.{frac_sec}]Z"
 * e.g. "2017-01-15T01:30:15.01Z"
 */
func DateRangeSearch(service *talent.Service, companyName string, startTime string, endTime string) (*talent.SearchJobsResponse, error) {
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
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with date range [%s, %s], Err: %v", startTime, endTime, err)
	}
	return resp, err
}

// [END date_range_filter]

// [START language_code_filter]

/**
 * Search on language code
 */
func LanguageCodeSearch(service *talent.Service, companyName string, languageCodes []string) (*talent.SearchJobsResponse, error) {
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
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with languange codes %v, Err: %v", languageCodes, err)
	}
	return resp, err
}

// [END language_code_filter]

// [START company_display_name_filter]

/**
 * Search on company display name
 */
func CompanyDisplayNameSearch(service *talent.Service, companyName string, companyDisplayNames []string) (*talent.SearchJobsResponse, error) {
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
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with company display names %v, Err: %v", companyDisplayNames, err)
	}
	return resp, err
}

// [END company_display_name_filter]

// [START compensation_fiter]

/**
 * Search on compensation
 */
func CompensationSearch(service *talent.Service, companyName string) (*talent.SearchJobsResponse, error) {
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
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with compensation, Err: %v", err)
	}
	return resp, err
}

// [END compensation_filter]

// [START general_search_sample_entry]

/**
 * Sample entry for the general job search
 */
func GeneralSearchSampleEntry() {
	service, _ := CreateCtsService()

	// Create a company before creating jobs
	companyToCreate := ConstructCompanyWithRequiredFields()
	companyCreated, _ := CreateCompany(service, companyToCreate)
	fmt.Printf("CreateCompany: %s\n", companyCreated.DisplayName)

	// Create a job
	jobTitle := "Systems Administrator"
	jobToCreate := ConstructJobWithRequiredFields(companyCreated.Name, jobTitle)
	jobToCreate.LanguageCode = "en-US"
	jobToCreate.EmploymentTypes = []string{"FULL_TIME"}
	compensationInfo := &talent.CompensationInfo{}
	compensationInfo.Entries = []*talent.CompensationEntry{
		&talent.CompensationEntry{
			Type: "BASE",
			Unit: "HOURLY",
			Amount: &talent.Money{
				CurrencyCode: "USD",
				Units:        12,
			},
		},
	}
	jobToCreate.CompensationInfo = compensationInfo

	jobCreated, _ := CreateJob(service, jobToCreate)
	fmt.Printf("CreateJob: %s\n", jobCreated.Title)

	// Wait for 10 seconds for post processing
	time.Sleep(10 * time.Second)

	resp, _ := BasicJobSearch(service, companyCreated.Name, jobTitle)
	fmt.Printf("BasicJobSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	categories := []string{"COMPUTER_AND_IT"}
	resp, _ = CategoryFilterSearch(service, companyCreated.Name, categories)
	fmt.Printf("CategoryFilterSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	employmentTypes := []string{"FULL_TIME", "CONTRACTOR", "PER_DIEM"}
	resp, _ = EmploymentTypesSearch(service, companyCreated.Name, employmentTypes)
	fmt.Printf("EmploymentTypesSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	resp, _ = DateRangeSearch(service, companyCreated.Name, "2000-01-01T00:00:00.01Z", "2099-01-01T00:00:00.01Z")
	fmt.Printf("DateRangeSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	resp, _ = LanguageCodeSearch(service, companyCreated.Name, []string{"pt-BR", "en-US"})
	fmt.Printf("LanguageCodeSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	resp, _ = CompanyDisplayNameSearch(service, companyCreated.Name, []string{"Google Sample"})
	fmt.Printf("CompanyDisplayNameSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	resp, _ = CompensationSearch(service, companyCreated.Name)
	fmt.Printf("CompensationSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	empty, _ := DeleteJob(service, jobCreated.Name)
	fmt.Printf("DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	empty, _ = DeleteCompany(service, companyCreated.Name)
	fmt.Printf("DeleteCompany StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)

}

// [END general_search_sample_entry]
