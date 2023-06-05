// Copyright 2023 Google LLC
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

package search

// [START genappbuilder_search]
import (
	"context"
	"fmt"

	discoveryengine "cloud.google.com/go/discoveryengine/apiv1beta"
	discoveryenginepb "cloud.google.com/go/discoveryengine/apiv1beta/discoveryenginepb"
	"google.golang.org/api/iterator"
)

// search searches for a query in a search engine given the Google Cloud Project ID,
// Location, and Search Engine ID.
//
// This example uses the default search engine.
func search(projectID, location, searchEngineID, query string) error {

	ctx := context.Background()

	// Create a client
	client, err := discoveryengine.NewSearchClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// Full resource name of search engine serving config
	servingConfig := fmt.Sprintf("projects/%s/locations/%s/collections/default_collection/dataStores/%s/servingConfigs/default_serving_config",
		projectID, location, searchEngineID)

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
