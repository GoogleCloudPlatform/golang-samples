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

// [START job_search_quick_start_list_companies]

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	talent "cloud.google.com/go/talent/apiv4beta1"
	"google.golang.org/api/iterator"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
)

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	parent := fmt.Sprintf("projects/%s", projectID)

	// Initialize job search client.
	ctx := context.Background()
	c, err := talent.NewCompanyClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	req := &talentpb.ListCompaniesRequest{
		Parent: parent,
	}
	it := c.ListCompanies(ctx, req)

	// Print the returned companies.
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			fmt.Printf("Done.\n")
			break
		}
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Printf("Company: %v\n", resp.GetName())
	}
}

// [END job_search_quick_start_list_companies]
