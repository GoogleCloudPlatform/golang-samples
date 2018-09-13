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

// [START basic_job]

// constructJobWithRequiredFields constructs a basic job with given companyName and jobTitle.
func constructJobWithRequiredFields(companyName string, jobTitle string) *talent.Job {
	requisitionID := fmt.Sprintf("sample-job-required-fields-%d", time.Now().UnixNano())
	applicationInfo := &talent.ApplicationInfo{
		Uris: []string{"https://googlesample.com/career"},
	}
	job := &talent.Job{
		RequisitionId:   requisitionID,
		Title:           jobTitle,
		CompanyName:     companyName,
		ApplicationInfo: applicationInfo,
		Description:     "Design, devolop, test, deploy, maintain and improve software.",
	}
	return job
}

// [END basic_job]

// [START create_job]

// createJob create a job as given.
func createJob(service *talent.Service, parent string, jobToCreate *talent.Job) (*talent.Job, error) {
	createJobRequest := &talent.CreateJobRequest{
		Job: jobToCreate,
	}
	job, err := service.Projects.Jobs.Create(parent, createJobRequest).Do()
	if err != nil {
		log.Fatalf("Failed to create job %q, Err: %v", jobToCreate.RequisitionId, err)
	}
	return job, err
}

// [END create_job]

// [START get_job]

// getJob gets a job by name.
func getJob(service *talent.Service, jobName string) (*talent.Job, error) {
	job, err := service.Projects.Jobs.Get(jobName).Do()
	if err != nil {
		log.Fatalf("Failed to get job %s, Err: %v", jobName, err)
	}

	return job, err
}

// [END get_job]

// [START update_job]

// updateJob update a job with all fields except name
func updateJob(service *talent.Service, jobName string, jobToUpdate *talent.Job) (*talent.Job, error) {
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

// updateJobWithMask updates a job by name with specific fields
// mask: comma separated top-level fields of Job
func updateJobWithMask(service *talent.Service, jobName string, mask string, jobToUpdate *talent.Job) (*talent.Job, error) {
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

// deleteJob deletes an existing job by name
func deleteJob(service *talent.Service, jobName string) (*talent.Empty, error) {
	empty, err := service.Projects.Jobs.Delete(jobName).Do()
	if err != nil {
		log.Fatalf("Failed to delete job %s, Err: %v", jobName, err)
	}

	return empty, err
}

// [END delete_job]

// [START list_jobs]

// listJobs lists jobs with companyName as filter
// filter required, eg companyName = "projects/api-test-project/companies/123"
func listJobs(service *talent.Service, parent string, filter string) (*talent.ListJobsResponse, error) {
	resp, err := service.Projects.Jobs.List(parent).Filter(filter).Do()
	if err != nil {
		log.Fatalf("Failed to list jobs with filter: %s, Err: %v", filter, err)
	}

	return resp, err
}

// [END list_jobs]

// [START run_basic_job_sample]

func runBasicJobSample(w io.Writer, projectID string) {
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

	// Construct a job
	jobTitle := "Software Engineer"
	jobToCreate := constructJobWithRequiredFields(companyCreated.Name, jobTitle)

	// Create a job
	jobCreated, err := createJob(service, parent, jobToCreate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "CreateJob: %s\n", jobCreated.Title)

	// Get an existing job
	jobName := jobCreated.Name
	jobGot, err := getJob(service, jobName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "GetJob: %s\n", jobGot.Title)

	// Update an existing job
	jobToUpdate := jobGot
	jobToUpdate.Title = "Software Engineer (updated)"
	jobUpdated, err := updateJob(service, jobName, jobToUpdate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "UpdateJob: %s\n", jobUpdated.Title)

	// Update job with field mask, only top level fields could be masked
	jobToUpdate.Title = "Software Engineer (unintended)"
	jobToUpdate.Department = "Engineering (updated with mask)"
	jobUpdatedWithMask, err := updateJobWithMask(service, jobName, "Department", jobToUpdate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "UpdateJobWithMask: Title: %s Department: %s\n", jobUpdatedWithMask.Title, jobUpdatedWithMask.Department)

	companyFilter := fmt.Sprintf("companyName = \"%s\"", companyCreated.Name)
	resp, err := listJobs(service, parent, companyFilter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "ListJobs Request ID: %q\n", resp.Metadata.RequestId)

	for _, job := range resp.Jobs {
		fmt.Fprintf(w, "-- Job: %q\n", job.Name)
	}

	// Delete Job
	empty, err := deleteJob(service, jobCreated.Name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "DeleteJob StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
	// Delete Company
	empty, err = deleteCompany(service, companyCreated.Name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "DeleteCompany StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)
}

// [END run_basic_job_sample]
