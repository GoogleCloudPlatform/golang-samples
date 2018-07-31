// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START quick_start]

// This is a quickstart sample of using the Google Cloud Job Discovery API.
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

func makeJob(companyName string) *jobs.Job {
	return &jobs.Job{
		RequisitionId:   fmt.Sprintf("job:%d", rand.New(rand.NewSource(time.Now().UnixNano())).Int63()),
		Title:        "Software Engineer",
		CompanyName:     companyName,
		ApplicationInfo: &jobs.ApplicationInfo{Uris: []string{"http://careers.google.com"}},
		Description:     "Design, develop, test, deploy, maintain and improve software.",
	}
}

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

// [START basic_keyword_search]

func basicKeywordSearch(ctx context.Context, js *jobs.Service, keyWord string, companyName string, parent string) (*jobs.SearchJobsResponse, error) {
	fmt.Println("Attempting to search jobs...")

	createCtx, cancel := context.WithTimeout(ctx, requestDeadline)
	defer cancel()

	requestMetadata := &jobs.RequestMetadata {
		UserId: "HashedUserId",
		SessionId: "HashedSessionId",
		Domain: "www.google.com",
	}

	query := &jobs.JobQuery {
		Query: keyWord,
	}
	if companyName != "" {
		query.CompanyNames = []string{companyName}
	}
	request := &jobs.SearchJobsRequest {
		SearchMode: "JOB_SEARCH",
		RequestMetadata: requestMetadata,
		JobQuery: query,
	}

	j, err := js.Projects.Jobs.Search(parent, request).Context(createCtx).Do()
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n %v\n", prettyFormat(j))
	return j, nil
}

// [END basic_keyword_search]

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

	jobToBeCreated := makeJob(company.Name)
	jobToBeCreated.Title = "Systems Administrator"
	job, err := createJob(ctx, jobService, jobToBeCreated, parent)
	if err != nil {
		fmt.Println("Failed to create a Job: ", err)
		return
	}

	time.Sleep(10 * time.Second)
	keyword := "Systems Administrator"
	basicKeywordSearch(ctx, jobService, keyword, company.Name, parent)

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

// [END quick_start]
