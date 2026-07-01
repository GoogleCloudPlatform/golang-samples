// Copyright 2026 Google LLC
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

package retail

import (
	"context"
	"fmt"

	retail "cloud.google.com/go/retail/apiv2"

	retailpb "cloud.google.com/go/retail/apiv2/retailpb"
)

// [START retail_v2_search_pagination]
// searchPagination method searches for products with pagination using Vertex AI Search for commerce.
// Performs a search request, then uses the nextToken to get the next page.
//
// projectID: The Google Cloud project ID.
// query: The search term for text search.
// visitorID: A unique identifier for the user.
func searchPagination(projectID, query, visitorID string) error {
	ctx := context.Background()

	client, err := retail.NewSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	placement := fmt.Sprintf("projects/%s/locations/global/catalogs/default_catalog/placements/default_search", projectID)

	// First page request
	firstPageReq := &retailpb.SearchRequest{
		Placement: placement,
		Query:     query,
		VisitorId: visitorID,
		PageSize:  5,
	}

	it := client.Search(ctx, firstPageReq)

	resp, err := it.Next()
	if err != nil {
		return fmt.Errorf("failed to fetch first page: %w", err)
	}

	fmt.Printf("First result from page 1: %s\n", resp.GetProduct().GetName())

	nextToken := it.PageInfo().Token
	if nextToken == "" {
		fmt.Println("No more pages available")
		return nil
	}
	fmt.Printf("Next page token: %s\n", nextToken)

	// Second page request using PageToken
	secondPageReq := &retailpb.SearchRequest{
		Placement: placement,
		Query:     query,
		VisitorId: visitorID,
		PageSize:  5,
		PageToken: nextToken,
	}

	it2 := client.Search(ctx, secondPageReq)
	resp2, err := it2.Next()
	if err != nil {
		return fmt.Errorf("failed to fetch first page: %w", err)
	}

	fmt.Printf("First result from page 2: %s\n", resp2.GetProduct().GetName())

	return nil
}

// [END retail_v2_search_pagination]
