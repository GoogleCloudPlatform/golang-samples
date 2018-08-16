package cjdsample

import (
	"fmt"
	"log"
	"time"

	talent "google.golang.org/api/jobs/v3"
)

// [START basic_job]

/**
 * Construct a basic job with given companyName and jobTitle
 */
func ConstructJobWithRequiredFields(companyName string, jobTitle string) *talent.Job {
	requisitionId := fmt.Sprintf("sample-job-required-fields-%d", time.Now().UnixNano())
	applicationInfo := &talent.ApplicationInfo{
		Uris: []string{"https://googlesample.com/career"},
	}
	job := &talent.Job{
		RequisitionId:   requisitionId,
		Title:           jobTitle,
		CompanyName:     companyName,
		ApplicationInfo: applicationInfo,
		Description:     "Design, devolop, test, deploy, maintain and improve software.",
	}
	//	fmt.Printf("Job constructed: %v\n",job)
	return job
}

// [END basic_job]

// [START create_job]

/**
 * Create a job
 */
func CreateJob(service *talent.Service, jobToCreate *talent.Job) (*talent.Job, error) {
	createJobRequest := &talent.CreateJobRequest{
		Job: jobToCreate,
	}
	job, err := service.Projects.Jobs.Create(GetParent(), createJobRequest).Do()
	if err != nil {
		log.Fatalf("Failed to create job %q, Err: %v", jobToCreate.RequisitionId, err)
	}
	return job, err
}

// [END create_job]

// [START get_job]

/**
 * Get a job
 */
func GetJob(service *talent.Service, jobName string) (*talent.Job, error) {
	job, err := service.Projects.Jobs.Get(jobName).Do()
	if err != nil {
		log.Fatalf("Failed to get job %s, Err: %v", jobName, err)
	}

	return job, err
}

// [END get_job]

// [START update_job]

/**
 * Update a job with all fields
 */
func UpdateJob(service *talent.Service, jobName string, jobToUpdate *talent.Job) (*talent.Job, error) {
	updateJobRequest := &talent.UpdateJobRequest{
		Job: jobToUpdate,
	}
	job, err := service.Projects.Jobs.Patch(jobName, updateJobRequest).Do()
	if err != nil {
		log.Fatalf("Failed to update job %s, Err: %v", jobName, err)
	}

	return job, err
}

// [END update_job]

// [START update_job_with_field_mask]

/**
 * Update a job
 * mask: comma separated top-level fields of Job
 */
// mask: comma separated top-level fields of Job
func UpdateJobWithMask(service *talent.Service, jobName string, mask string, jobToUpdate *talent.Job) (*talent.Job, error) {
	updateJobRequest := &talent.UpdateJobRequest{
		Job:        jobToUpdate,
		UpdateMask: mask,
	}
	job, err := service.Projects.Jobs.Patch(jobName, updateJobRequest).Do()
	if err != nil {
		log.Fatalf("Failed to update job %s with field mask %s, Err: %v", jobName, mask, err)
	}

	return job, err
}

// [END update_job_with_field_mask]

// [START delete_job]

/**
 * Delete a job
 */
func DeleteJob(service *talent.Service, jobName string) (*talent.Empty, error) {
	empty, err := service.Projects.Jobs.Delete(jobName).Do()
	if err != nil {
		log.Fatalf("Failed to delete job %s, Err: %v", jobName, err)
	}

	return empty, err
}

// [END delete_job]

// [START list_jobs]

/**
 * List jobs with companyName as filter
 * filter required, eg companyName = "projects/api-test-project/companies/123"
 */
func ListJobs(service *talent.Service, filter string) (*talent.ListJobsResponse, error) {
	resp, err := service.Projects.Jobs.List(GetParent()).Filter(filter).Do()
	if err != nil {
		log.Fatalf("Failed to list jobs with filter: %s, Err: %v", filter, err)
	}

	return resp, err
}

// [END list_jobs]

// [START basic_job_sample_entry]

func BasicJobSampleEntry() {
	service, _ := CreateCtsService()

	// Create a company before creating jobs
	companyToCreate := ConstructCompanyWithRequiredFields()
	companyCreated, _ := CreateCompany(service, companyToCreate)
	fmt.Printf("CreateCompany: %s\n", companyCreated.DisplayName)

	// Construct a job
	jobTitle := "Software Engineer"
	jobToCreate := ConstructJobWithRequiredFields(companyCreated.Name, jobTitle)

	// Create a job
	jobCreated, _ := CreateJob(service, jobToCreate)
	fmt.Printf("CreateJob: %s\n", jobCreated.Title)

	// Get an existing job
	jobName := jobCreated.Name
	jobGot, _ := GetJob(service, jobName)
	fmt.Printf("GetJob: %s\n", jobGot.Title)

	// Update an existing job
	jobToUpdate := jobGot
	jobToUpdate.Title = "Software Engineer (updated)"
	jobUpdated, _ := UpdateJob(service, jobName, jobToUpdate)
	fmt.Printf("UpdateJob: %s\n", jobUpdated.Title)

	// Update job with field mask, only top level fields could be masked
	jobToUpdate.Title = "Software Engineer (unintended)"
	jobToUpdate.Department = "Engineering (updated with mask)"
	jobUpdatedWithMask, _ := UpdateJobWithMask(service, jobName, "Department", jobToUpdate)
	fmt.Printf("UpdateJobWithMask: Title: %s Department: %s\n", jobUpdatedWithMask.Title, jobUpdatedWithMask.Department)

	companyFilter := fmt.Sprintf("companyName = \"%s\"", companyCreated.Name)
	resp, _ := ListJobs(service, companyFilter)
	fmt.Printf("ListJobs Request ID: %q\n", resp.Metadata.RequestId)

	for _, job := range resp.Jobs {
		fmt.Printf("-- Job: %q\n", job.Name)
	}

	// Delete Job
	empty, _ := DeleteJob(service, jobCreated.Name)
	fmt.Printf("DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	// Delete Company
	empty, _ = DeleteCompany(service, companyCreated.Name)
	fmt.Printf("DeleteCompany StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
}

// [END basic_job_sample_entry]
