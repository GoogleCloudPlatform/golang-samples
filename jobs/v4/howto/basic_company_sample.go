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

	// Initializes a companyService client.
	c, err := talent.NewCompanyClient(ctx)
	if err != nil {
		fmt.Printf("talent.NewCompanyClient: %v", err)
		return nil, err
	}

  // Construct a createCompany request.
	req := &talentpb.CreateCompanyRequest{
		Parent: "projects/" + projectId,
		Company: companyToCreate,
	}

	resp, err := c.CreateCompany(ctx, req)
	if err != nil {
		fmt.Printf("talent.NewCompanyClient: %v", err)
		return nil, err
	}

	fmt.Printf("Created company: %q\n", resp.GetName())

	return resp, nil
}

// [END job_search_create_company]

// [START job_search_get_company]

// getCompany gets an existing company by name.
func getCompany(w io.Writer, companyName string) (*talentpb.Company, error) {
	ctx := context.Background()

	// Initialize a companyService client.
	c, err := talent.NewCompanyClient(ctx)
	if err != nil {
		fmt.Printf("talent.NewCompanyClient: %v", err)
		return nil, err
	}

	// Construct a getCompany request.
	req := &talentpb.GetCompanyRequest{
		// The resource name of the company to be retrieved.
    // The format is "projects/{project_id}/companies/{company_id}".
		Name: companyName,
	}

	resp, err := c.GetCompany(ctx, req)
	if err != nil {
		fmt.Printf("failed to get company %q: %v", companyName, err)
		return nil, err
	}

	fmt.Fprintf(w, "Company: %q\n", resp.GetName())

	return resp, nil
}

// [END job_search_get_company]

// [START job_search_delete_company]

// deleteCompany deletes an existing company by name. Companies with
// existing jobs cannot be deleted until those jobs have been deleted.
func deleteCompany(w io.Writer, companyName string) error {
	ctx := context.Background()

	// Initialize a companyService client.
	c, err := talent.NewCompanyClient(ctx)
	if err != nil {
		fmt.Printf("talent.NewCompanyClient: %v", err)
		return err
	}

	// Construct a deleteCompany request.
	req := &talentpb.DeleteCompanyRequest{
		// The resource name of the company to be deleted.
		// The format is "projects/{project_id}/companies/{company_id}".
		Name: companyName,
	}

	if err := c.DeleteCompany(ctx, req); err != nil {
		fmt.Printf("failed to delete company %q: %v", companyName, err)
		return err
	}

	fmt.Printf("Deleted company: %q\n", companyName)

	return nil
}

// [END job_search_delete_company

// [START job_search_list_companies]

// listCompanies lists all companies in the project.
func listCompanies(w io.Writer, projectId string) error {
	ctx := context.Background()

	// Initialize a compnayService client.
	c, err := talent.NewCompanyClient(ctx)
	if err != nil {
		fmt.Printf("talent.NewCompanyClient: %v", err)
		return err
	}

	// Construct a listCompanies request.
	req := &talentpb.ListCompaniesRequest{
		Parent: "projects/" + projectId,
	}

	it := c.ListCompanies(ctx, req)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			fmt.Printf("it.Next: %q", err)
			return err
		}
		fmt.Fprintf(w, "Listed company: %q\n", resp.GetName())
	}
}

// [END job_search_list_companies]
