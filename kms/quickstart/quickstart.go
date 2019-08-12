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

// [START kms_quickstart]

// Sample quickstart is a basic program that uses Cloud KMS.
package main

import (
	"context"
	"fmt"
	"log"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"google.golang.org/api/iterator"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

func main() {
	projectID := "your-project-id"
	// Location of the key rings.
	locationID := "global"

	// Create the KMS client.
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// The resource name of the key rings.
	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)

	// Build the request.
	req := &kmspb.ListKeyRingsRequest{
		Parent: parent,
	}
	// Query the API.
	it := client.ListKeyRings(ctx, req)

	// Iterate and print results.
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to list key rings: %v", err)
		}
		fmt.Printf("KeyRing: %q\n", resp.Name)
	}
}

// [END kms_quickstart]
