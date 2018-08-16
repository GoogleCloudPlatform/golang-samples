package cjdsample

import (
	"fmt"
	"log"
	"time"

	talent "google.golang.org/api/jobs/v3"
)

// [START featured_job]

/**
 * Construct a job as featured
 */
func ConstructFeaturedJob(companyName string, jobTitle string) (*talent.Job) {
	requisitionId := fmt.Sprintf("featured-job-required-fields-%d", time.Now().UnixNano())
	applicationInfo := &talent.ApplicationInfo {
		Uris: []string {"https://googlesample.com/career"},
	}
	job := &talent.Job{
		RequisitionId: requisitionId,
		Title: jobTitle,
		CompanyName: companyName,
		ApplicationInfo: applicationInfo,
		Description: "Design, devolop, test, deploy, maintain and improve software.",
		PromotionValue: 2,
	}
//	fmt.Printf("Job constructed: %v\n",job)
	return job
}
// [END featured_job]

// [START search_featured_job]

/**
  * Simple search jobs with keyword.
  */
func SearchFeaturedJobs(service *talent.Service, companyName string, query string) (*talent.SearchJobsResponse, error) {
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
		Query: query,
	}
	if companyName != "" {
		jobQuery.CompanyNames = []string{companyName}
	}

	searchJobsRequest := &talent.SearchJobsRequest {
		RequestMetadata: requestMetadata,
		// Set the actual search term as defined in the jobQurey
		JobQuery: jobQuery, 
		// Set the search mode to a featured search
		// which would only search the jobs with positive promotion value
		SearchMode: "FEATURED_JOB_SEARCH",
	}
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with query %v, Err: %v", query, err)
	}
	return resp, err
}
// [END search_featured_job]


func FeaturedJobSearchSampleEntry() {
	service, _ := CreateCtsService()

	// Create a company before creating jobs
	companyToCreate := ConstructCompanyWithRequiredFields()
	companyCreated, _ := CreateCompany(service, companyToCreate)
	fmt.Printf("CreateCompany: %s\n", companyCreated.DisplayName)

	// Construct a job
	jobTitle := "Software Engineer (Featured)"
	jobToCreate := ConstructFeaturedJob(companyCreated.Name, jobTitle)

	// Create a featured job
	jobCreated, _ := CreateJob(service, jobToCreate)
	fmt.Printf("CreateJob: %s\n", jobCreated.Title)

	// Wait for 10 seconds for post processing
	time.Sleep(10 * time.Second)

	// Search for featured jobs
	resp, _ := SearchFeaturedJobs(service, companyCreated.Name, jobTitle)
	fmt.Printf("SearchFeaturedJobs StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	// Delete Job
	empty, _ := DeleteJob(service, jobCreated.Name)
	fmt.Printf("DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	// Delete Company
	emptyResp, _ := DeleteCompany(service, companyCreated.Name)
	fmt.Printf("DeleteCompany StatusCode: %d\n", emptyResp.ServerResponse.HTTPStatusCode)
}




