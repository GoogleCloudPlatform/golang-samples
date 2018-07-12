// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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

// [START basic_company]

func makeCompany() *jobs.Company {
	return &jobs.Company{
		DisplayName:          "Google",
		HqLocation:           "1600 Amphitheatre Parkway Mountain View, CA 94043",
		DistributorCompanyId: fmt.Sprintf("company:%d", rand.New(rand.NewSource(time.Now().UnixNano())).Int63()),
	}
}

// [END basic_company]

// [START create_company]

func createCompany(ctx context.Context, js *jobs.Service, toCreate *jobs.Company) (*jobs.Company, error) {
	fmt.Println("Attempting to create a Company...")

	createCtx, cancel := context.WithTimeout(ctx, requestDeadline)
	defer cancel()

	c, err := js.Companies.Create(toCreate).Context(createCtx).Do()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Company created:\n %v\n", prettyFormat(c))
	return c, nil
}

// [END create_company]

// [START get_company]

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

// [END get_company]

// [START update_company]

func patchCompany(ctx context.Context, js *jobs.Service, company *jobs.Company) (*jobs.Company, error) {
	return patchCompanyWithFieldMask(ctx, js, company, "")
}

// [END update_company]

// [START update_company_with_field_mask]

func patchCompanyWithFieldMask(ctx context.Context, js *jobs.Service, company *jobs.Company, fields string) (*jobs.Company, error) {
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

	fmt.Printf("Company updated:\n %v\n", prettyFormat(c))
	return c, nil
}

// [END update_company_with_field_mask]

// [START delete_company]
func deleteCompany(ctx context.Context, js *jobs.Service, name string) error {
	fmt.Printf("Attempting to delete a Company with name %s...\n", name)

	deleteCtx, cancel := context.WithTimeout(ctx, requestDeadline)
	defer cancel()

	if _, err := js.Companies.Delete(name).Context(deleteCtx).Do(); err != nil {
		return err
	}

	fmt.Println("Company deleted")
	return nil
}

// [END delete_company]

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

	company, err := createCompany(ctx, jobService, makeCompany())
	if err != nil {
		fmt.Println("Failed to create company: ", err)
		return
	}

	company, err = getCompany(ctx, jobService, company.Name)
	if err != nil {
		fmt.Printf("Failed to get a company with name %s: %v\n", company.Name, err)
		return
	}

	company.Website = "https://elgoog.im/"
	company, err = patchCompany(ctx, jobService, company)
	if err != nil {
		fmt.Printf("Failed to update a company with name %s: %v\n", company.Name, err)
		return
	}

	company.Website = "https://google.com/"
	company.Suspended = true
	// Only the website field should be updated by this call.
	company, err = patchCompanyWithFieldMask(ctx, jobService, company, "website")
	if err != nil {
		fmt.Printf("Failed to update a company with name %s: %v\n", company.Name, err)
		return
	}

	if err := deleteCompany(ctx, jobService, company.Name); err != nil {
		fmt.Printf("Failed to delete a company with name %s: %v\n", company.Name, err)
		return
	}
}
