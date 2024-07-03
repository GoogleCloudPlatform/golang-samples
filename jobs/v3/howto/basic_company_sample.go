// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package howto

import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	talent "google.golang.org/api/jobs/v3"
)

// [START job_create_company]

// createCompany creates a company as given.
func createCompany(w io.Writer, projectID string, companyToCreate *talent.Company) (*talent.Company, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %w", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %w", err)
	}

	parent := "projects/" + projectID
	req := &talent.CreateCompanyRequest{
		Company: companyToCreate,
	}
	company, err := service.Projects.Companies.Create(parent, req).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create company %q: %w", companyToCreate.DisplayName, err)
	}

	return company, nil
}

// [END job_create_company]

// [START job_get_company]

// getCompany gets an existing company by name.
func getCompany(w io.Writer, name string) (*talent.Company, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %w", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %w", err)
	}

	company, err := service.Projects.Companies.Get(name).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get company %q: %w", name, err)
	}

	fmt.Fprintf(w, "Company: %q\n", company.Name)

	return company, nil
}

// [END job_get_company]

// [START job_update_company]

// updateCompany update a company with all fields.
func updateCompany(w io.Writer, name string, companyToUpdate *talent.Company) (*talent.Company, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %w", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %w", err)
	}

	updateCompanyRequest := &talent.UpdateCompanyRequest{
		Company: companyToUpdate,
	}
	company, err := service.Projects.Companies.Patch(name, updateCompanyRequest).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to update company %q: %w", name, err)
	}

	return company, nil
}

// [END job_update_company]

// [START job_update_company_with_field_mask]

// updateCompanyWithMask updates a company with specific fields.
// mask is a comma separated list of top-level fields of talent.Company.
func updateCompanyWithMask(w io.Writer, name string, mask string, companyToUpdate *talent.Company) (*talent.Company, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %w", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %w", err)
	}

	req := &talent.UpdateCompanyRequest{
		Company:    companyToUpdate,
		UpdateMask: mask,
	}
	company, err := service.Projects.Companies.Patch(name, req).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to update company %q with mask %q: %w", name, mask, err)
	}

	return company, nil
}

// [END job_update_company_with_field_mask]

// [START job_delete_company]

// deleteCompany deletes an existing company by name.
func deleteCompany(w io.Writer, name string) error {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return fmt.Errorf("google.DefaultClient: %w", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return fmt.Errorf("talent.New: %w", err)
	}

	if _, err := service.Projects.Companies.Delete(name).Do(); err != nil {
		return fmt.Errorf("failed to delete company %q: %w", name, err)
	}

	return nil
}

// [END job_delete_company]

// [START job_list_companies]

// listCompanies lists all companies in the project.
func listCompanies(w io.Writer, projectID string) (*talent.ListCompaniesResponse, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %w", err)
	}
	// Create the jobs service client.
	service, err := talent.New(client)
	if err != nil {
		return nil, fmt.Errorf("talent.New: %w", err)
	}

	parent := "projects/" + projectID
	resp, err := service.Projects.Companies.List(parent).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list companies: %w", err)
	}

	fmt.Fprintln(w, "Companies:")
	for _, c := range resp.Companies {
		fmt.Fprintf(w, "\t%q\n", c.Name)
	}

	return resp, nil
}

// [END job_list_companies]
