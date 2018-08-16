package cjdsample

import (
	"fmt"
	"log"
	"time"

	googleapi "google.golang.org/api/googleapi"
	talent "google.golang.org/api/jobs/v3"
)

func ConstructJobWithCustomAttributes(companyName string, jobTitle string) (*talent.Job) {
	// requisition id shoud be the unique Id in your system
	requisitionId := fmt.Sprintf("job-with-custom-attribute-%d", time.Now().UnixNano())
	applicationInfo := &talent.ApplicationInfo {
		Uris: []string {"https://googlesample.com/career"},
	}
	customAttrStr := &talent.CustomAttribute {
		Filterable: true,
		StringValues: []string{"someStrVal"},
	}
	customAttrLong := &talent.CustomAttribute {
		Filterable: true,
		LongValues: googleapi.Int64s {900},
	}

	customAttributes := map[string]talent.CustomAttribute {
		"someFieldString": *customAttrStr,
		"someFieldLong": *customAttrLong,
	}

	job := &talent.Job{
		RequisitionId: requisitionId,
		Title: jobTitle,
		CompanyName: companyName,
		ApplicationInfo: applicationInfo,
		Description: "Design, devolop, test, deploy, maintain and improve software.",
		CustomAttributes: customAttributes,
	}
//	fmt.Printf("Job constructed: %v\n",job)
	return job
}


func FilterOnStringValueCustomAttribute(service *talent.Service)(*talent.SearchJobsResponse, error) {
	requestMetadata := &talent.RequestMetadata {
		UserId: "HashedUsrId",
		SessionId: "HashedSessionId",
		Domain: "www.googlesample.com",
	}

	customAttrFilter := "NOT EMPTY(someFieldString)"
	query := &talent.JobQuery {
		CustomAttributeFilter: customAttrFilter,
	}
	searchJobsRequest := &talent.SearchJobsRequest {
		JobQuery: query,
		RequestMetadata: requestMetadata,
		JobView: "JOB_VIEW_FULL",
	}
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with string value custom attribute, Err: %v", err)
	}

	return resp, err
}


func FilterOnLongValueCustomAttribute(service *talent.Service)(*talent.SearchJobsResponse, error) {
	requestMetadata := &talent.RequestMetadata {
		UserId: "HashedUsrId",
		SessionId: "HashedSessionId",
		Domain: "www.googlesample.com",
	}

	customAttrFilter := "someFieldLong < 1000"
	query := &talent.JobQuery {
		CustomAttributeFilter: customAttrFilter,
	}
	searchJobsRequest := &talent.SearchJobsRequest {
		JobQuery: query,
		RequestMetadata: requestMetadata,
		JobView: "JOB_VIEW_FULL",
	}
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with long value custom attribute, Err: %v", err)
	}

	return resp, err
}

func FilterOnMultiCustomAttributes(service *talent.Service) (*talent.SearchJobsResponse, error) {
	requestMetadata := &talent.RequestMetadata {
		UserId: "HashedUsrId",
		SessionId: "HashedSessionId",
		Domain: "www.googlesample.com",
	}

	customAttrFilter := "(someFieldString = \"someStrVal\") AND (someFieldLong < 1000)"
	query := &talent.JobQuery {
		CustomAttributeFilter: customAttrFilter,
	}

	searchJobsRequest := &talent.SearchJobsRequest {
		JobQuery: query,
		RequestMetadata: requestMetadata,
		JobView: "JOB_VIEW_FULL",
	}
	resp, err := service.Projects.Jobs.Search(GetParent(), searchJobsRequest).Do()
	if err != nil {
		log.Fatalf("Failed to search for jobs with multiple custom attributes, Err: %v", err)
	}

	return resp, err
}


func CustomAttributeSampleEntry() {
	service, _ := CreateCtsService()

	// Create a company before creating jobs
	companyToCreate := ConstructCompanyWithRequiredFields()
	companyCreated, _ := CreateCompany(service, companyToCreate)
	fmt.Printf("CreateCompany: %s\n", companyCreated.DisplayName)

	// Create a job with custom fields
	jobTitle := "Software Engineer"
	jobToCreate := ConstructJobWithCustomAttributes(companyCreated.Name, jobTitle)
	jobCreated, _ := CreateJob(service, jobToCreate)
	fmt.Printf("CreateJob: %s\n", jobCreated.Title)

	time.Sleep(10 * time.Second)

	resp, _ := FilterOnStringValueCustomAttribute(service)
	fmt.Printf("FilterOnStringValueCustomAttribute StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	resp, _ = FilterOnLongValueCustomAttribute(service)
	fmt.Printf("FilterOnLongValueCustomAttribute StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	resp, _ = FilterOnMultiCustomAttributes(service)
	fmt.Printf("FilterOnMultiCustomAttributes StatusCode: %d\n", resp.ServerResponse.HTTPStatusCode)
	fmt.Printf("MatchingJobs size: %d\n", len(resp.MatchingJobs))
	for _, mJob := range resp.MatchingJobs {
		fmt.Printf("-- match job: %s\n", mJob.Job.Title)
	}

	// Delete Job
	empty, _ := DeleteJob(service, jobCreated.Name)
	fmt.Printf("DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	// Delete Company
	empty, _ = DeleteCompany(service, companyCreated.Name)
	fmt.Printf("DeleteCompany StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)

}
