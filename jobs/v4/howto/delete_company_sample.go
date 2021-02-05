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

// [START job_search_delete_company]
import (
	"context"
	"fmt"
	"io"

	talent "cloud.google.com/go/talent/apiv4beta1"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
)

// deleteCompany deletes an existing company. Companies with
// existing jobs cannot be deleted until those jobs have been deleted.
func deleteCompany(w io.Writer, projectID, companyID string) error {
	ctx := context.Background()

	// Initialize a companyService client.
	c, err := talent.NewCompanyClient(ctx)
	if err != nil {
		return fmt.Errorf("talent.NewCompanyClient: %v", err)
	}

	// Construct a deleteCompany request.
	companyName := fmt.Sprintf("projects/%s/companies/%s", projectID, companyID)
	req := &talentpb.DeleteCompanyRequest{
		// The resource name of the company to be deleted.
		// The format is "projects/{project_id}/companies/{company_id}".
		Name: companyName,
	}

	if err := c.DeleteCompany(ctx, req); err != nil {
		return fmt.Errorf("DeleteCompany(%s): %v", companyName, err)
	}

	fmt.Fprintf(w, "Deleted company: %q\n", companyName)

	return nil
}

// [END job_search_delete_company
