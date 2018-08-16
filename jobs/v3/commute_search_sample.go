package cjdsample

import (
	"fmt"
	"log"
	"time"

	talent "google.golang.org/api/jobs/v3"
)

// [START commute_search]

/**
 * Search on commute search
 */
func CommuteSearch(service *talent.Service, companyName string) (*talent.SearchJobsResponse, error) {
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
		CommuteFilter: &talent.CommuteFilter {
			RoadTraffic: "TRAFFIC_FREE",
			CommuteMethod: "TRANSIT",
			TravelDuration: "1000s",
			StartCoordinates: &talent.LatLng {
				Latitude: 37.422408,
				Longitude: -122.085609,
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
		RequirePreciseResultSize: true,
	}
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with commute filter, Err: %v", err)
	}
	return resp, err
}

// [END commute_search]

func CommuteSearchSampleEntry() {
	service, _ := CreateCtsService()

	// Create a company before creating jobs
	companyToCreate := ConstructCompanyWithRequiredFields()
	companyCreated, _ := CreateCompany(service, companyToCreate)
	fmt.Printf("CreateCompany: %s\n", companyCreated.DisplayName)

	jobTitle := "Software Engineer"
	jobToCreate := ConstructJobWithRequiredFields(companyCreated.Name, jobTitle)
	jobToCreate.Addresses = []string{"1600 Amphitheatre Parkway, Mountain View, CA 94043"}
	jobCreated, _ := CreateJob(service, jobToCreate)
	fmt.Printf("CreateJob: %s\n", jobCreated.Title)

	// Wait several seconds for post processing
	time.Sleep(10 * time.Second)

	resp, _ := CommuteSearch(service, companyCreated.Name)
	fmt.Printf("CommuteSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	empty, _ := DeleteJob(service, jobCreated.Name)
	fmt.Printf("DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	empty, _ = DeleteCompany(service, companyCreated.Name)
	fmt.Printf("DeleteCompany StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
}
