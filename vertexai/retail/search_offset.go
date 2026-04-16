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
	"google.golang.org/api/iterator"
)

// [START retail_v2_search_offset]
// searchOffset method searches for products with an offset using Vertex AI Search for commerce.
//
// Performs a search request starting from a specified position.
//
// projectID: The Google Cloud project ID.
// query: The search term.
// visitorID: A unique identifier for the user.
// offset: The number of results to skip.for products using Vertex AI Search for commerce.
func searchOffset(projectID, query, visitorID string, offset int32) error {
	ctx := context.Background()

	client, err := retail.NewSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	req := &retailpb.SearchRequest{
		Placement: fmt.Sprintf("projects/%s/locations/global/catalogs/default_catalog/placements/default_search", projectID),
		Query:     query,
		VisitorId: visitorID,
		PageSize:  10,
		Offset:    offset,
	}

	it := client.Search(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("it.Next: %w", err)
		}
		fmt.Printf("Search item: %v\n", resp)
	}
	return nil
}

// [END retail_v2_search_offset]
