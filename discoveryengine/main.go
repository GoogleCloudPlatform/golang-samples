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

package main

// [START genappbuilder_search]
import (
	"context"
	"flag"
	"fmt"
	"log"

	genappbuilder "cloud.google.com/go/discoveryengine/apiv1beta"
	genappbuilderpb "cloud.google.com/go/discoveryengine/apiv1beta/discoveryenginepb"
	"google.golang.org/api/iterator"
)

func main() {
	projectID := flag.String("project", "YOUR_PROJECT_ID", "Google Cloud Project ID")
	location := flag.String("location", "global", "search engine location")
	searchEngineID := flag.String("searchengine", "YOUR_SEARCH_ENGINE_ID", "search engine ID")
	query := flag.String("query", "YOUR_SEARCH_QUERY", "search query")
	flag.Parse()

	ctx := context.Background()

	// Create a client
	client, err := genappbuilder.NewSearchClient(ctx)
	if err != nil {
		log.Fatalf("unable to create a search client: %v", err)
	}
	defer client.Close()

	// Full resource name of search engine serving config
	servingConfig := fmt.Sprintf("projects/%s/locations/%s/collections/default_collection/dataStores/%s/servingConfigs/default_serving_config", *projectID, *location, *searchEngineID)

	searchRequest := &genappbuilderpb.SearchRequest{
		ServingConfig: servingConfig,
		Query:         *query,
	}

	it := client.Search(ctx, searchRequest)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("unable to retrieve: %v", err)
		}
		fmt.Printf("%+v\n", resp)
	}
}
// [END genappbuilder_search]
