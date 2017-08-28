// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample quickstart is a basic program that uses Cloud KMS.
package main

// [START kms_quickstart]
import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

func main() {
	projectID := "your-project-id"
	// Location of the key rings.
	locationID := "global"

	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	// Create the KMS client.
	kmsService, err := cloudkms.New(client)
	if err != nil {
		log.Fatal(err)
	}

	// The resource name of the key rings.
	parentName := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)

	// Make the RPC call.
	response, err := kmsService.Projects.Locations.KeyRings.List(parentName).Do()
	if err != nil {
		log.Fatalf("Failed to list key rings: %v", err)
	}

	// Print the returned key rings.
	for _, keyRing := range response.KeyRings {
		fmt.Printf("KeyRing: %q\n", keyRing.Name)
	}
}

// [END kms_quickstart]
