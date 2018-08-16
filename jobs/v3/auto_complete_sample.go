package cjdsample


import (
	"fmt"
	"log"
	"time"

	talent "google.golang.org/api/jobs/v3"
)
// Auto completes job titles within given companyName
//
func JobTitleAutoComplete(service *talent.Service, companyName string, query string) (*talent.CompleteQueryResponse, error) {
	complete := service.Projects.Complete(GetParent()).Query(query).LanguageCode("en-US").Type("JOB_TITLE").PageSize(10)
	if companyName != "" {
		complete.CompanyName(companyName)
	}
	resp, err := complete.Do()
	if err != nil {
		log.Fatalf("Failed to auto complete with query %s in company %s, Err: %v", query, companyName, err)
	}
	return resp, err

}


// Auto completes job titles within given companyName
//
func DefaultAutoComplete(service *talent.Service, companyName string, query string) (*talent.CompleteQueryResponse, error) {
	complete := service.Projects.Complete(GetParent()).Query(query).LanguageCode("en-US").Type("COMBINED").PageSize(10)
	if companyName != "" {
		complete.CompanyName(companyName)
	}
	resp, err := complete.Do()
	if err != nil {
		log.Fatalf("Failed to auto complete with query %s in company %s, Err: %v", query, companyName, err)
	}
	return resp, err

}

func AutoCompleteSampleEntry() {
	service, _ := CreateCtsService()

	// Create a company before creating jobs
	companyToCreate := ConstructCompanyWithRequiredFields()
	companyCreated, _ := CreateCompany(service, companyToCreate)
	fmt.Printf("CreateCompany: %s\n", companyCreated.DisplayName)

	// Create a SDE job
	jobTitleSWE := "Software Engineer"
	jobToCreateSWE := ConstructJobWithRequiredFields(companyCreated.Name, jobTitleSWE)
	jobCreatedSWE, _ := CreateJob(service, jobToCreateSWE)
	fmt.Printf("CreateJob: %s\n", jobCreatedSWE.Title)
	// Create a PM job
	jobTitlePM := "GAP Product Manager"
	jobToCreatePM := ConstructJobWithRequiredFields(companyCreated.Name, jobTitlePM)
	jobCreatedPM, _ := CreateJob(service, jobToCreatePM)
	fmt.Printf("CreateJob: %s\n", jobCreatedPM.Title)

	// Wait several seconds for post processing
	time.Sleep(10 * time.Second)

	query := "sof"
	resp, _ := DefaultAutoComplete(service, "", query)
	fmt.Printf("DefaultAutoComplete query: %s StatusCode: %d\n", query, resp.ServerResponse.HTTPStatusCode)
	for _, comp := range resp.CompletionResults {
		fmt.Printf("-- suggestion: %s\n", comp.Suggestion)
	}

	query = "sof"
	resp, _ = JobTitleAutoComplete(service, "", query)
	fmt.Printf("JobTitleAutoComplete query: %s StatusCode: %d\n", query, resp.ServerResponse.HTTPStatusCode)
	for _, comp := range resp.CompletionResults {
		fmt.Printf("-- suggestion: %s\n", comp.Suggestion)
	}

	query = "gap"
	resp, _ = DefaultAutoComplete(service, companyCreated.Name, query)
	fmt.Printf("DefaultAutoComplete query: %s StatusCode: %d\n", query, resp.ServerResponse.HTTPStatusCode)
	for _, comp := range resp.CompletionResults {
		fmt.Printf("-- suggestion: %s\n", comp.Suggestion)
	}

	query = "gap"
	resp, _ = JobTitleAutoComplete(service, companyCreated.Name, query)
	fmt.Printf("JobTitleAutoComplete query: %s StatusCode: %d\n", query, resp.ServerResponse.HTTPStatusCode)
	for _, comp := range resp.CompletionResults {
		fmt.Printf("-- suggestion: %s\n", comp.Suggestion)
	}

	// Delete Job
	empty, _ := DeleteJob(service, jobCreatedSWE.Name)
	fmt.Printf("DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	// Delete Job
	empty, _ = DeleteJob(service, jobCreatedPM.Name)
	fmt.Printf("DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	// Delete Company
	empty, _ = DeleteCompany(service, companyCreated.Name)
	fmt.Printf("DeleteCompany StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)

}
