// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package discoveryengine

// [START genappbuilder_search]
import (
	"context"
	"fmt"

	discoveryengine "cloud.google.com/go/discoveryengine/apiv1"
	discoveryenginepb "cloud.google.com/go/discoveryengine/apiv1/discoveryenginepb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// search searches for a query in a search app given the Google Cloud Project ID,
// Location, and Data Store ID.
func search(projectID, location, dataStoreID, query string) error {

	ctx := context.Background()

	// Create a client
	endpoint := "discoveryengine.googleapis.com:443" // Default to global endpoint
	if location != "global" {
		endpoint = fmt.Sprintf("%s-%s", location, endpoint)
	}
	client, err := discoveryengine.NewSearchClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		fmt.Println(fmt.Errorf("error creating Vertex AI Search client: %w", err))
	}
	defer client.Close()

	// Full resource name of search engine serving config
	servingConfig := fmt.Sprintf("projects/%s/locations/%s/collections/default_collection/dataStores/%s/servingConfigs/default_config", projectID, location, dataStoreID)

	searchRequest := &discoveryenginepb.SearchRequest{
		ServingConfig: servingConfig,
		Query:         query,
	}

	it := client.Search(ctx, searchRequest)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", resp)
	}

	return nil
}

// [END genappbuilder_search]
