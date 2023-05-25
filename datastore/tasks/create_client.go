// Copyright 2023 Google LLC
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

// [START datastore_build_service]
import (
	"context"
	"log"

	"cloud.google.com/go/datastore"
)

func createClient(projectID string) (*datastore.Client, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}
	// Note: call the following from main() to ensure the client
	// properly frees all resources.
	// defer client.Close()
	return client, nil
}

// [END datastore_build_service]
