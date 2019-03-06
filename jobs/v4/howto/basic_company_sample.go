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

	talent "cloud.google.com/go/talent/apiv4beta1"
	"google.golang.org/api/iterator"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
)

// [START job_search_create_company]

// createCompany creates a company as given.
func createCompany(w io.Writer, projectId string, companyToCreate *talentpb.Company) (*talentpb.Company, error) {
	ctx := context.Background()

	// Create a company service client.
	c, err := talent.NewCompanyClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("talent.NewCompanyClient: %v", err)
	}

  // projectId := "Your Google Cloud Project ID"
	req := &talentpb.CreateCompanyRequest{
		Parent: "projects/" + projectId,
		Company: companyToCreate,
	}
	resp, err := c.CreateCompany(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Failed to create company: %v", err)
	}

	fmt.Printf("Creating company: %v\n", resp.GetDisplayName())
	fmt.Printf("Created company name: %v\n", resp)
	return resp, nil
}

// [END job_search_create_company]

// [START job_search_get_company]

// getCompany gets an existing company by name.
func getCompany(w io.Writer, name string) (*talentpb.Company, error) {
	ctx := context.Background()

	// Create a company service client.
	c, err := talent.NewCompanyClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("talent.NewCompanyClient: %v", err)
	}

	req := &talentpb.GetCompanyRequest{
		// The resource name of the company to be retrieved.
    // The format is "projects/{project_id}/companies/{company_id}".
		Name: name,
	}

	resp, err := c.GetCompany(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get company %q: %v", name, err)
	}

  fmt.Printf("\nGetting company: %v\n", resp.GetDisplayName())
	fmt.Fprintf(w, "Got company name: %q\n", resp.GetName())
	return resp, nil
}

// [END job_search_get_company]

// [START job_search_delete_company]

// deleteCompany deletes an existing company by name.
func deleteCompany(w io.Writer, name string) error {
	ctx := context.Background()

	// Create a company service client.
	c, err := talent.NewCompanyClient(ctx)
	if err != nil {
		return fmt.Errorf("talent.NewCompanyClient: %v", err)
	}

	req := &talentpb.DeleteCompanyRequest{
		Name: name,
	}
	if err := c.DeleteCompany(ctx, req); err != nil {
		return fmt.Errorf("failed to delete company %q: %v", name, err)
	}

	fmt.Printf("Deleted company: %s\n", name)
	return nil
}

// [END job_search_delete_company

// [START job_search_list_companies]

// listCompanies lists all companies in the project.
func listCompanies(w io.Writer, projectId string) error {
	ctx := context.Background()

	// Create a compnay service client.
	c, err := talent.NewCompanyClient(ctx)
	if err != nil {
		return fmt.Errorf("talent.NewCompanyClient: %v", err)
	}

	req := &talentpb.ListCompaniesRequest{
		Parent: fmt.Sprintf("projects/%s", projectId),
	}
	it := c.ListCompanies(ctx, req)
	// Print the returned companies.
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return fmt.Errorf("it.Next: %v", err)
		}
		fmt.Printf("\nListing company: %v\n%v\n", resp.GetDisplayName(), resp.GetName())
		fmt.Fprintf(w, "Listed company display name: %v\n", resp.GetName())
	}
}

// [END job_search_list_companies]
