// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
	parentName := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)

	// Build the request.
	req := &kmspb.ListKeyRingsRequest{
		Parent: parentName,
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
