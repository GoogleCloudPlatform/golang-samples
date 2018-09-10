// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	talent "google.golang.org/api/jobs/v3"
)

// [START create_service]

// createCtsService creates service of Cloud Talent Solution.
func createCtsService() (*talent.Service, error) {
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		log.Fatalf("Failed to create client of Cloud Talent Solution, Err: %v", err)
	}
	// Create the jobs service client.
	ctsService, err := talent.New(client)
	if err != nil {
		log.Fatalf("Failed to create service of Cloud Talent Solution, Err: %v", err)
	}
	return ctsService, err
}

// [END create_service]

// [START basic_company]

// constructCompanyWithRequiredFields constructs a company with required fields: ExternalId and DisplayName.
func constructCompanyWithRequiredFields() *talent.Company {
	externalId := fmt.Sprintf("sample-company-%d", time.Now().UnixNano())
	return &talent.Company{
		ExternalId:  externalId,
		DisplayName: "Google Sample",
	}
}

// [END basic_company]

// [START create_company]

// createCompany creates a company as given.
func createCompany(service *talent.Service, parent string, companyToCreate *talent.Company) (*talent.Company, error) {
	createCompanyRquest := &talent.CreateCompanyRequest{
		Company: companyToCreate,
	}
	company, err := service.Projects.Companies.Create(parent, createCompanyRquest).Do()
	if err != nil {
		log.Fatalf("Failed to create company %s, Err: %v", companyToCreate.DisplayName, err)
	}

	return company, err
}

// [END create_company]

// [START get_company]

// getCompany gets an existing company by name.
func getCompany(service *talent.Service, name string) (*talent.Company, error) {
	company, err := service.Projects.Companies.Get(name).Do()
	if err != nil {
		log.Fatalf("Failed to get company %s, Err: %v", name, err)
	}

	return company, err
}

// [END get_company]

// [START update_company]

// updateCompany update a company with all fields.
func updateCompany(service *talent.Service, name string, companyToUpdate *talent.Company) (*talent.Company, error) {
	updateCompanyRequest := &talent.UpdateCompanyRequest{
		Company: companyToUpdate,
	}
	company, err := service.Projects.Companies.Patch(name, updateCompanyRequest).Do()
	if err != nil {
		log.Fatalf("Failed to update company %s, Err: %v", name, err)
	}

	return company, err
}

// [END update_company]

// [START update_company_with_field_mask]

// updateCompanyWithMask updates a company with specific fields.
// mask: comma separated top-level fields of Company
func updateCompanyWithMask(service *talent.Service, name string, mask string, companyToUpdate *talent.Company) (*talent.Company, error) {
	updateCompanyRequest := &talent.UpdateCompanyRequest{
		Company:    companyToUpdate,
		UpdateMask: mask,
	}
	company, err := service.Projects.Companies.Patch(name, updateCompanyRequest).Do()
	if err != nil {
		log.Fatalf("Failed to update company %s with mask %s, Err: %v", name, mask, err)
	}

	return company, err
}

// [END update_company_with_field_mask]

// [START delete_company]

// deleteCompany deletes an existing company by name.
func deleteCompany(service *talent.Service, name string) (*talent.Empty, error) {
	empty, err := service.Projects.Companies.Delete(name).Do()
	if err != nil {
		log.Fatalf("Failed to delete company %s, Err: %v", name, err)
	}

	return empty, err
}

// [END delete_company

// [START list_companies]

// listCompanies lists all companies in the project
func listCompanies(service *talent.Service, parent string) (*talent.ListCompaniesResponse, error) {
	resp, err := service.Projects.Companies.List(parent).Do()
	if err != nil {
		log.Fatalf("Failed to list companies, Err: %v", err)
	}

	return resp, err
}

// [END list_companies]

// [START run_basic_company_sample]

func runBasicCompanySample(w io.Writer) {
	parent := fmt.Sprintf("projects/%s", os.Getenv("GOOGLE_CLOUD_PROJECT"))
	service, _ := createCtsService()

	companyToCreate := constructCompanyWithRequiredFields()
	companyCreated, _ := createCompany(service, parent, companyToCreate)
	fmt.Fprintf(w, "CreateCompany: %s\n", companyCreated.DisplayName)

	name := companyCreated.Name
	companyGot, _ := getCompany(service, name)
	fmt.Fprintf(w, "GetCompany: %s\n", companyGot.DisplayName)

	companyToUpdate := companyCreated
	companyToUpdate.DisplayName = "Google Sample (updated)"
	companyUpdated, _ := updateCompany(service, name, companyToUpdate)
	fmt.Fprintf(w, "UpdateCompany: %s\n", companyUpdated.DisplayName)

	companyUpdated.WebsiteUri = "http://googlesample.com"
	companyUpdated.DisplayName = "Google Sample (updated with mask)"
	companyUpdatedWithMask, _ := updateCompanyWithMask(service, name, "WebSiteUri,DisplayName", companyUpdated)
	fmt.Fprintf(w, "UpdateCompanyWithMask: %s\n", companyUpdatedWithMask.DisplayName)

	empty, _ := deleteCompany(service, name)
	fmt.Fprintf(w, "DeleteCompany StatusCode: %d\n", empty.ServerResponse.HTTPStatusCode)

	resp, _ := listCompanies(service, parent)
	fmt.Fprintf(w, "ListCompanies Request ID: %q\n", resp.Metadata.RequestId)

	for _, company := range resp.Companies {
		fmt.Fprintf(w, "-- Company: %q\n", company.Name)
	}

}

// [END run_basic_company_sample]
