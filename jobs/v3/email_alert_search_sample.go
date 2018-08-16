package cjdsample

import (
	"fmt"
	"log"
	"time"

	talent "google.golang.org/api/jobs/v3"
)

// [START search_for_alerts]

/**
 * Search jobs for alert
 */
func SearchForAlerts(service *talent.Service, companyName string) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata {
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
	}

	searchJobsRequest := &talent.SearchJobsRequest {
		RequestMetadata: requestMetadata,
		// Set the search mode to a regular search
		SearchMode: "JOB_SEARCH",
	}
	if companyName != "" {
		jobQuery := &talent.JobQuery {
			CompanyNames: []string{companyName},
		}
		searchJobsRequest.JobQuery = jobQuery
	}

	resp, err := service.Projects.Jobs.SearchForAlert(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with alerts, Err: %v", err)
	}
	return resp, err
}

// [END search_for_alerts]

func EmailAlertSearchSampleEntry() {
	service, _ := CreateCtsService()

	// Create a company before creating jobs
	companyToCreate := ConstructCompanyWithRequiredFields()
	companyCreated, _ := CreateCompany(service, companyToCreate)
	fmt.Printf("CreateCompany: %s\n", companyCreated.DisplayName)

	// Create a SDE job
	jobTitle := "Software Engineer"
	jobToCreate := ConstructJobWithRequiredFields(companyCreated.Name, jobTitle)
	jobCreated, _ := CreateJob(service, jobToCreate)
	fmt.Printf("CreateJob: %s\n", jobCreated.Title)

	// Wait for 10 seconds for post processing
	time.Sleep(10 * time.Second)

	// Search jobs with alerts
	resp, _ := SearchForAlerts(service, companyCreated.Name)
	fmt.Printf("SearchForAlerts StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	empty, _ := DeleteJob(service, jobCreated.Name)
	fmt.Printf("DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	empty, _ = DeleteCompany(service, companyCreated.Name)
	fmt.Printf("DeleteCompany StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)

}

