// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START basic-company]
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/oauth2/google"
	jobs "google.golang.org/api/jobs/v2"
)

const requestDeadline = 10 * time.Second

func prettyFormat(v interface{}) string {
	j, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	return string(j)
}

func createCompany(ctx context.Context, js *jobs.Service) (*jobs.Company, error) {
	fmt.Println("Attempting to create a Company...")

	createCtx, cancel := context.WithTimeout(ctx, requestDeadline)
	defer cancel()

	c, err := js.Companies.Create(&jobs.Company{
		DisplayName:          "Google",
		HqLocation:           "1600 Amphitheatre Parkway Mountain View, CA 94043",
		DistributorCompanyId: fmt.Sprintf("company:%d", rand.New(rand.NewSource(time.Now().UnixNano())).Int63()),
	}).Context(createCtx).Do()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Company created:\n %v\n", prettyFormat(c))
	return c, nil
}

func getCompany(ctx context.Context, js *jobs.Service, name string) (*jobs.Company, error) {
	fmt.Printf("Attempting to get a Company with name %s...\n", name)

	getCtx, cancel := context.WithTimeout(ctx, requestDeadline)
	defer cancel()

	c, err := js.Companies.Get(name).Context(getCtx).Do()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Company retrieved:\n %v\n", prettyFormat(c))
	return c, nil
}

func patchCompany(ctx context.Context, js *jobs.Service, company *jobs.Company, fields string) (*jobs.Company, error) {
	fmt.Printf("Attempting to patch a Company with name %s...\n", company.Name)

	patchCtx, cancel := context.WithTimeout(ctx, requestDeadline)
	defer cancel()

	req := js.Companies.Patch(company.Name, company).Context(patchCtx)
	if fields != "" {
		req.UpdateCompanyFields(fields)
	}

	c, err := req.Do()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Company after patching:\n %v\n", prettyFormat(c))
	return c, nil
}

func main() {
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, jobs.CloudPlatformScope)
	if err != nil {
		fmt.Println("Failed to create an authenticated HTTP client: ", err)
		return
	}

	// Create the jobs service client.
	jobService, err := jobs.New(client)
	if err != nil {
		fmt.Println("Failed to create a JobService client: ", err)
		return
	}

	company, err := createCompany(ctx, jobService)
	if err != nil {
		fmt.Println("Failed to create company: ", err)
		return
	}

	company, err = getCompany(ctx, jobService, company.Name)
	if err != nil {
		fmt.Printf("Failed to get a company with name %s: \n", company.Name, err)
		return
	}

	company.Website = "https://elgoog.im/"
	company, err = patchCompany(ctx, jobService, company, "")
	if err != nil {
		fmt.Printf("Failed to get a company with name %s: \n", company.Name, err)
		return
	}

	company.Website = "https://google.com/"
	company.Suspended = true
	// Only the website field should be updated by this call.
	company, err = patchCompany(ctx, jobService, company, "website")
	if err != nil {
		fmt.Printf("Failed to get a company with name %s: \n", company.Name, err)
		return
	}
}

// [END basic-company]
