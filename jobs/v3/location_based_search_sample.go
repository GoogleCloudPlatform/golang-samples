package cjdsample

import (
	"fmt"
	"log"
	"time"

	talent "google.golang.org/api/jobs/v3"
)

// [START basic_location_search]

/**
 * Basic location search
 */
func BasicLocationSearch(service *talent.Service, companyName string, location string, distance float64) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata {
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
	}
	jobQuery := &talent.JobQuery {
		LocationFilters: []*talent.LocationFilter { 
			&talent.LocationFilter {
				Address: location,
				DistanceInMiles: distance,
			},
		},
	}

	if companyName != "" {
		jobQuery.CompanyNames = []string{companyName}
	}

	searchJobsRequest := &talent.SearchJobsRequest {
		RequestMetadata: requestMetadata,
		// Set the actual search term as defined in the jobQurey
		JobQuery: jobQuery, 
		// Set the search mode to a regular search
		SearchMode: "JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with basic location %s within %f miles, Err: %v", location, distance, err)
	}
	return resp, err
}

// [END basic_location_search]

// [START city_location_search]

/**
 * City location search
 */
func CityLocationSearch(service *talent.Service, companyName string, location string) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata {
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
	}
	jobQuery := &talent.JobQuery {
		LocationFilters: []*talent.LocationFilter { 
			&talent.LocationFilter {
				Address: location,
			},
		},
	}

	if companyName != "" {
		jobQuery.CompanyNames = []string{companyName}
	}

	searchJobsRequest := &talent.SearchJobsRequest {
		RequestMetadata: requestMetadata,
		// Set the actual search term as defined in the jobQurey
		JobQuery: jobQuery, 
		// Set the search mode to a regular search
		SearchMode: "JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with city location %s, Err: %v", location, err)
	}
	return resp, err
}

// [END city_location_search]

// [START broadening_location_search]

/**
 * Broadening location search
 */
func BroadeningLocationSearch(service *talent.Service, companyName string, location string) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata {
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
	}
	jobQuery := &talent.JobQuery {
		LocationFilters: []*talent.LocationFilter { 
			&talent.LocationFilter {
				Address: location,
			},
		},
	}

	if companyName != "" {
		jobQuery.CompanyNames = []string{companyName}
	}

	searchJobsRequest := &talent.SearchJobsRequest {
		RequestMetadata: requestMetadata,
		// Set the actual search term as defined in the jobQurey
		JobQuery: jobQuery, 
		// Set the search mode to a regular search
		SearchMode: "JOB_SEARCH",
		EnableBroadening: true,
	}
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with broadening location %v, Err: %v", location, err)
	}
	return resp, err
}

// [END broadening_location_search]

// [START keyword_location_search]

/**
 * Keyword location search
 */
func KeywordLocationSearch(service *talent.Service, companyName string, location string, distance float64, keyword string) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata {
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
	}
	jobQuery := &talent.JobQuery {
		LocationFilters: []*talent.LocationFilter { 
			&talent.LocationFilter {
				Address: location,
				DistanceInMiles: distance,
			},
		},
		Query: keyword,
	}

	if companyName != "" {
		jobQuery.CompanyNames = []string{companyName}
	}

	searchJobsRequest := &talent.SearchJobsRequest {
		RequestMetadata: requestMetadata,
		// Set the actual search term as defined in the jobQurey
		JobQuery: jobQuery, 
		// Set the search mode to a regular search
		SearchMode: "JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with keyword %s in location %v within %f miles, Err: %v", keyword, location, distance, err)
	}
	return resp, err
}

// [END keyword_location_search]

// [START multi_locations_search]

/**
 * Multi locations search
 */
func MultiLocationsSearch(service *talent.Service, companyName string, location string, distance float64, location2 string) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata {
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
	}
	jobQuery := &talent.JobQuery {
		LocationFilters: []*talent.LocationFilter { 
			&talent.LocationFilter {
				Address: location,
				DistanceInMiles: distance,
			},
			&talent.LocationFilter {
				Address: location,
			},
		},
	}

	if companyName != "" {
		jobQuery.CompanyNames = []string{companyName}
	}

	searchJobsRequest := &talent.SearchJobsRequest {
		RequestMetadata: requestMetadata,
		// Set the actual search term as defined in the jobQurey
		JobQuery: jobQuery, 
		// Set the search mode to a regular search
		SearchMode: "JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with multi locations %s and %s within %f miles, Err: %v", location, location2, distance, err)
	}
	return resp, err
}

// [END keyword_locations_search]

func LocationBasedSearchSampleEntry() {
	service, _ := CreateCtsService()

	// Create a company before creating jobs
	companyToCreate := ConstructCompanyWithRequiredFields()
	companyCreated, _ := CreateCompany(service, companyToCreate)
	fmt.Printf("CreateCompany: %s\n", companyCreated.DisplayName)

	location := "Mountain View, CA"
	distance := 0.5
	keyword := "Software Engineer"
	location2 := "Sunnyvale, CA"

	// Create a SDE job
	jobTitle := keyword
	jobToCreate := ConstructJobWithRequiredFields(companyCreated.Name, jobTitle)
	jobToCreate.Addresses = []string {location}

	jobCreated, _ := CreateJob(service, jobToCreate)
	fmt.Printf("CreateJob: %s\n", jobCreated.Title)

	jobTitle2 := "Senior " + keyword
	jobToCreate2 := ConstructJobWithRequiredFields(companyCreated.Name, jobTitle2)
	jobToCreate2.Addresses = []string {location2}
	jobCreated2, _ := CreateJob(service, jobToCreate2)
	fmt.Printf("CreateJob: %s\n", jobCreated2.Title)

	// Wait for 10 seconds for post processing
	time.Sleep(10 * time.Second)

	resp, _ := BasicLocationSearch(service, companyCreated.Name, location, distance)
	fmt.Printf("BasicLocationSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	resp, _ = CityLocationSearch(service, companyCreated.Name, location)
	fmt.Printf("CityLocationSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	resp, _ = BroadeningLocationSearch(service, companyCreated.Name, location)
	fmt.Printf("BroadeningLocationSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	resp, _ = KeywordLocationSearch(service, companyCreated.Name, location, distance, keyword)
	fmt.Printf("KeywordLocationSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	resp, _ = MultiLocationsSearch(service, companyCreated.Name, location, distance, location2)
	fmt.Printf("MultiLocationsSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	empty, _ := DeleteJob(service, jobCreated.Name)
	fmt.Printf("DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	empty, _ = DeleteJob(service, jobCreated2.Name)
	fmt.Printf("DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	empty, _ = DeleteCompany(service, companyCreated.Name)
	fmt.Printf("DeleteCompany StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)

}
