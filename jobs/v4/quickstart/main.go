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

package main

// [START job_search_quick_start_list_companies]
import (
	"context"
	"fmt"
	"os"

	talent "cloud.google.com/go/talent/apiv4beta1"
	"cloud.google.com/go/talent/apiv4beta1/talentpb"
	"google.golang.org/api/iterator"
)

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	// Initialize job search client.
	ctx := context.Background()
	c, err := talent.NewCompanyClient(ctx)
	if err != nil {
		fmt.Printf("talent.NewCompanyClient: %v", err)
		return
	}
	defer c.Close()

	// Construct a listCompany request.
	req := &talentpb.ListCompaniesRequest{
		Parent: "projects/" + projectID,
	}

	it := c.ListCompanies(ctx, req)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			fmt.Printf("Done.\n")
			break
		}
		if err != nil {
			fmt.Printf("it.Next: %v", err)
			return
		}
		fmt.Printf("Company: %v\n", resp.GetName())
	}
}

// [END job_search_quick_start_list_companies]
