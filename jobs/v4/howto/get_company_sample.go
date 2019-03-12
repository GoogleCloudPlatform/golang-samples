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

// [START job_search_get_company]
import (
	"context"
	"fmt"
	"io"

	talent "cloud.google.com/go/talent/apiv4beta1"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
)

// getCompany gets an existing company by its resource name.
func getCompany(w io.Writer, projectID, companyID string) (*talentpb.Company, error) {
	ctx := context.Background()

	// Initialize a companyService client.
	c, err := talent.NewCompanyClient(ctx)
	if err != nil {
		fmt.Printf("talent.NewCompanyClient: %v\n", err)
		return nil, err
	}

	// Construct a getCompany request.
	companyName := fmt.Sprintf("projects/%s/companies/%s", projectID, companyID)
	req := &talentpb.GetCompanyRequest{
		// The resource name of the company to be retrieved.
		// The format is "projects/{project_id}/companies/{company_id}".
		Name: companyName,
	}

	resp, err := c.GetCompany(ctx, req)
	if err != nil {
		fmt.Printf("failed to get company %q: %v\n", companyName, err)
		return nil, err
	}

	fmt.Fprintf(w, "Company: %q\n", resp.GetName())
	fmt.Fprintf(w, "Company display name: %q\n", resp.GetDisplayName())

	return resp, nil
}

// [END job_search_get_company]
