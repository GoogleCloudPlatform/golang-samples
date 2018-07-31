// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START quick_start]

// This is a quickstart sample of using the Google Cloud Job Discovery API.
package main

import (
	"os"
	"fmt"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	talent "google.golang.org/api/jobs/v3"
)

func main() {
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	// Create the jobs service client.
	cloudTalentSolutionSerivce, err := talent.New(client)
	if err != nil {
		log.Fatal(err)
	}

	parent := fmt.Sprintf("projects/%s", projectID)

	// Make the RPC call.
	response, err := cloudTalentSolutionSerivce.Projects.Companies.List(parent).Do()
	if err != nil {
		log.Fatalf("Failed to list Companies: %v", err)
	}

    // Print the request id.
	fmt.Printf("Request ID: %q\n", response.Metadata.RequestId)

	// Print the returned companies.
	for _, company := range response.Companies {
		fmt.Printf("Company: %q\n", company.Name)
	}
}

// [END quick_start]
