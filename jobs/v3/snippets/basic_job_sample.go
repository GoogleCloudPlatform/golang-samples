// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/oauth2/google"
	jobs "google.golang.org/api/jobs/v3"
)

const requestDeadline = 10 * time.Second

func prettyFormat(v interface{}) string {
	j, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	return string(j)
}

func makeCompany() *jobs.Company {
	return &jobs.Company{
		DisplayName:          "Google",
		HeadquartersAddress:           "1600 Amphitheatre Parkway Mountain View, CA 94043",
		ExternalId: fmt.Sprintf("company:%d", rand.New(rand.NewSource(time.Now().UnixNano())).Int63()),
	}
}


func createCompany(ctx context.Context, js *jobs.Service, toCreate *jobs.Company, parent string) (*jobs.Company, error) {
	fmt.Println("Attempting to create a Company...")

	createCtx, cancel := context.WithTimeout(ctx, requestDeadline)
	defer cancel()

	request := &jobs.CreateCompanyRequest {
		Company: toCreate,
	}

	c, err := js.Projects.Companies.Create(parent,request).Context(createCtx).Do()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Company created:\n %v\n", prettyFormat(c))
	return c, nil
}


func deleteCompany(ctx context.Context, js *jobs.Service, name string) error {
	fmt.Printf("Attempting to delete a Company with name %s...\n", name)

	deleteCtx, cancel := context.WithTimeout(ctx, requestDeadline)
	defer cancel()

	if _, err := js.Projects.Companies.Delete(name).Context(deleteCtx).Do(); err != nil {
		return err
	}

	fmt.Println("Company deleted")
	return nil
}

// [START basic_job]

func makeJob(companyName string) *jobs.Job {
	return &jobs.Job{
		RequisitionId:   fmt.Sprintf("job:%d", rand.New(rand.NewSource(time.Now().UnixNano())).Int63()),
		Title:        "Software Engineer",
		CompanyName:     companyName,
		ApplicationInfo: &jobs.ApplicationInfo{Uris: []string{"http://careers.google.com"}},
		Description:     "Design, develop, test, deploy, maintain and improve software.",
	}
}

// [END basic_job]

// [START create_job]

func createJob(ctx context.Context, js *jobs.Service, toCreate *jobs.Job, parent string) (*jobs.Job, error) {
	fmt.Println("Attempting to create a Job...")

	createCtx, cancel := context.WithTimeout(ctx, requestDeadline)
	defer cancel()

	j, err := js.Projects.Jobs.Create(parent, &jobs.CreateJobRequest{Job: toCreate}).Context(createCtx).Do()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Job created:\n %v\n", prettyFormat(j))
	return j, nil
}

// [END create_job]

// [START get_job]

func getJob(ctx context.Context, js *jobs.Service, name string) (*jobs.Job, error) {
	fmt.Printf("Attempting to get a Job with name %s...\n", name)

	getCtx, cancel := context.WithTimeout(ctx, requestDeadline)
	defer cancel()

	j, err := js.Projects.Jobs.Get(name).Context(getCtx).Do()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Job retrieved:\n %v\n", prettyFormat(j))
	return j, nil
}

// [END get_job]

// [START update_job]

func patchJob(ctx context.Context, js *jobs.Service, job *jobs.Job) (*jobs.Job, error) {
	return patchJobWithFieldMask(ctx, js, job, "")
}

// [END update_job]

// [START update_job_with_field_mask]

func patchJobWithFieldMask(ctx context.Context, js *jobs.Service, job *jobs.Job, fields string) (*jobs.Job, error) {
	fmt.Printf("Attempting to patch a Job with name %s...\n", job.Name)

	patchCtx, cancel := context.WithTimeout(ctx, requestDeadline)
	defer cancel()

	req := &jobs.UpdateJobRequest{Job: job}
	if fields != "" {
		req.UpdateMask = fields
	}

	call := js.Projects.Jobs.Patch(job.Name, req).Context(patchCtx)
	j, err := call.Do()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Job updated:\n %v\n", prettyFormat(j))
	return j, nil
}

// [END update_job_with_field_mask]

// [START delete_job]
func deleteJob(ctx context.Context, js *jobs.Service, name string) error {
	fmt.Printf("Attempting to delete a Job with name %s...\n", name)

	deleteCtx, cancel := context.WithTimeout(ctx, requestDeadline)
	defer cancel()

	if _, err := js.Projects.Jobs.Delete(name).Context(deleteCtx).Do(); err != nil {
		return err
	}

	fmt.Println("Job deleted")
	return nil
}

// [END delete_job]

func main() {
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, jobs.CloudPlatformScope)
	if err != nil {
		fmt.Println("Failed to create an authenticated HTTP client: ", err)
		return
	}

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	parent := fmt.Sprintf("projects/%s", projectID)

	// Create the jobs service client.
	jobService, err := jobs.New(client)
	if err != nil {
		fmt.Println("Failed to create a JobService client: ", err)
		return
	}

	company, err := createCompany(ctx, jobService, makeCompany(), parent)
	if err != nil {
		fmt.Println("Failed to create a Company: ", err)
		return
	}

	job, err := createJob(ctx, jobService, makeJob(company.Name), parent)
	if err != nil {
		fmt.Println("Failed to create a Job: ", err)
		return
	}

	job, err = getJob(ctx, jobService, job.Name)
	if err != nil {
		fmt.Printf("Failed to get a job with name %s: %v\n", job.Name, err)
		return
	}

	job.Description = "changedDescription"
	job, err = patchJob(ctx, jobService, job)
	if err != nil {
		fmt.Printf("Failed to update a job with name %s: %v\n", job.Name, err)
		return
	}

	// Only the job title field should be updated by this call.
	job.CompanyName = "changedCompanyName"
	job.Title = "changedJobTitle"
	job, err = patchJobWithFieldMask(ctx, jobService, job, "jobTitle")
	if err != nil {
		fmt.Printf("Failed to update a job with name %s: %v\n", job.Name, err)
		return
	}

	if err := deleteJob(ctx, jobService, job.Name); err != nil {
		fmt.Printf("Failed to delete a job with name %s: %v\n", job.Name, err)
		return
	}

	// Delete company only after cleaning up all jobs under this company.
	if err := deleteCompany(ctx, jobService, company.Name); err != nil {
		fmt.Printf("Failed to delete a company with name %s: %v\n", company.Name, err)
		return
	}
}