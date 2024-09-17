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

// [START genappbuilder_create_data_store]
import (
	"context"
	"fmt"

	discoveryengine "cloud.google.com/go/discoveryengine/apiv1"
	discoveryenginepb "cloud.google.com/go/discoveryengine/apiv1/discoveryenginepb"
	"google.golang.org/api/option"
)

// search searches for a query in a search app given the Google Cloud Project ID,
// Location, and App ID.
func createDataStore(projectID, location, dataStoreID string) error {

	ctx := context.Background()

	// Create a client
	endpoint := "discoveryengine.googleapis.com:443" // Default to global endpoint
	if location != "global" {
		endpoint = fmt.Sprintf("%s-%s", location, endpoint)
	}
	client, err := discoveryengine.NewDataStoreClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		fmt.Println(fmt.Errorf("error creating Vertex AI Search client: %w", err))
	}
	defer client.Close()

	// The full resource name of the collection.
	// e.g. projects/project/locations/{location}/collections/default_collection
	parent := fmt.Sprintf("projects/%s/locations/%s/collections/default_collection", projectID, location)

	dataStore := &discoveryenginepb.DataStore{
		DisplayName: "My Data Store",
		// Options: GENERIC, MEDIA, HEALTHCARE_FHIR
		IndustryVertical: discoveryenginepb.IndustryVertical_GENERIC,
		// Options: SOLUTION_TYPE_RECOMMENDATION, SOLUTION_TYPE_SEARCH, SOLUTION_TYPE_CHAT, SOLUTION_TYPE_GENERATIVE_CHAT
		SolutionTypes: []discoveryenginepb.SolutionType{discoveryenginepb.SolutionType_SOLUTION_TYPE_SEARCH},
		// TODO(developer): Update content_config based on data store type.
		// Options: NO_CONTENT, CONTENT_REQUIRED, PUBLIC_WEBSITE
		ContentConfig: discoveryenginepb.DataStore_CONTENT_REQUIRED,
	}

	request := &discoveryenginepb.CreateDataStoreRequest{
		Parent:      parent,
		DataStoreId: dataStoreID,
		DataStore:   dataStore,
	}

	// Make the request.
	op, err := client.CreateDataStore(ctx, request)
	if err != nil {
		return fmt.Errorf("client.CreateDataStore: %w", err)
	}

	fmt.Printf("Waiting for operation to complete: %s\n", op.Name())
	// After the operation is complete, get information from operation metadata.
	_, err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("op.Wait: %w", err)
	}
	return nil
}

// [END genappbuilder_create_data_store]
