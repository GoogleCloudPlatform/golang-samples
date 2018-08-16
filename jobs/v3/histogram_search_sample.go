package cjdsample

import (
	"fmt"
	"log"
	"time"

	talent "google.golang.org/api/jobs/v3"
)

// [START commute_search]

/**
 * Search on histogram search
 */
func HistogramSearch(service *talent.Service, companyName string) (*talent.SearchJobsResponse, error) {
	// Make sure to set the requestMetadata the same as the associated search request
	requestMetadata := &talent.RequestMetadata {
		// Make sure to hash your userID
		UserId: "HashedUsrId",
		// Make sure to hash the sessionID
		SessionId: "HashedSessionId",
		// Domain of the website where the search is conducted
		Domain: "www.googlesample.com",
	}

	histogramFacets := &talent.HistogramFacets {
		SimpleHistogramFacets: []string{"COMPANY_ID"},
		CustomAttributeHistogramFacets: []*talent.CustomAttributeHistogramRequest {
			&talent.CustomAttributeHistogramRequest {
				Key: "someFieldString",
				StringValueHistogram: true,
			},
		},
	}

	searchJobsRequest := &talent.SearchJobsRequest {
		RequestMetadata: requestMetadata,
		HistogramFacets: histogramFacets, 
		// Set the search mode to a regular search
		SearchMode: "JOB_SEARCH",
		RequirePreciseResultSize: true,
	}
	if companyName != "" {
		jobQuery := &talent.JobQuery {
			CompanyNames: []string{companyName},
		}
		searchJobsRequest.JobQuery = jobQuery
	}

	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with Historgram Facets, Err: %v", err)
	}
	return resp, err
}

// [END histogram_search]

func HistogramSearchSampleEntry() {
	service, _ := CreateCtsService()

	// Create a company before creating jobs
	companyToCreate := ConstructCompanyWithRequiredFields()
	companyCreated, _ := CreateCompany(service, companyToCreate)
	fmt.Printf("CreateCompany: %s\n", companyCreated.DisplayName)

	// Create a SDE job
	jobTitleSWE := "Software Engineer"
	jobToCreateSWE := ConstructJobWithCustomAttributes(companyCreated.Name, jobTitleSWE)
	jobCreatedSWE, _ := CreateJob(service, jobToCreateSWE)
	fmt.Printf("CreateJob: %s\n", jobCreatedSWE.Title)

	// Wait several seconds for post processing
	time.Sleep(10 * time.Second)

	resp, _ := HistogramSearch(service, companyCreated.Name)
	fmt.Printf("HistogramSearch StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}
	fmt.Printf("SimpleHistogramResults size: %d\n", len(resp.HistogramResults.SimpleHistogramResults))
	for _, hist := range resp.HistogramResults.SimpleHistogramResults {
		fmt.Printf("-- simple histogram searchType: %s value: %v\n", hist.SearchType, hist.Values) 
	}
	fmt.Printf("CustomAttributeHistogramResults size: %d\n", len(resp.HistogramResults.CustomAttributeHistogramResults))
	for _, hist := range resp.HistogramResults.CustomAttributeHistogramResults {
		fmt.Printf("-- custom-attribute histogram key: %s value: %v\n", hist.Key, hist.StringValueHistogramResult) 
	}


	empty, _ := DeleteJob(service, jobCreatedSWE.Name)
	fmt.Printf("DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	empty, _ = DeleteCompany(service, companyCreated.Name)
	fmt.Printf("DeleteCompany StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
}
