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

// [START job_search_create_company]
import (
	"context"
	"fmt"
	"io"

	talent "cloud.google.com/go/talent/apiv4beta1"
	"cloud.google.com/go/talent/apiv4beta1/talentpb"
)

// createCompany creates a company as given.
func createCompany(w io.Writer, projectID, externalID, displayName string) (*talentpb.Company, error) {
	ctx := context.Background()

	// Initializes a companyService client.
	c, err := talent.NewCompanyClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("talent.NewCompanyClient: %w", err)
	}
	defer c.Close()

	// Construct a createCompany request.
	req := &talentpb.CreateCompanyRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		Company: &talentpb.Company{
			ExternalId:  externalID,
			DisplayName: displayName,
		},
	}

	resp, err := c.CreateCompany(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("CreateCompany: %w", err)
	}

	fmt.Fprintf(w, "Created company: %q\n", resp.GetName())

	return resp, nil
}

// [END job_search_create_company]
