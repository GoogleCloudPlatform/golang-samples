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

// [START job_search_list_company]
import (
	"context"
	"fmt"
	"io"

	talent "cloud.google.com/go/talent/apiv4beta1"
	"google.golang.org/api/iterator"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
)

// listCompanies lists all companies in the project.
func listCompanies(w io.Writer, projectID string) error {
	ctx := context.Background()

	// Initialize a compnayService client.
	c, err := talent.NewCompanyClient(ctx)
	if err != nil {
		fmt.Printf("talent.NewCompanyClient: %v\n", err)
		return err
	}

	// Construct a listCompanies request.
	req := &talentpb.ListCompaniesRequest{
		Parent: "projects/" + projectID,
	}

	it := c.ListCompanies(ctx, req)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			fmt.Printf("it.Next: %q\n", err)
			return err
		}
		fmt.Fprintf(w, "Listing company: %q\n", resp.GetName())
		fmt.Fprintf(w, "Display name: %v\n", resp.GetDisplayName())
		fmt.Fprintf(w, "External ID: %v\n", resp.GetExternalId())
	}
}

// [END job_search_list_companies]
