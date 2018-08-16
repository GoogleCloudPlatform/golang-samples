// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START quick_start]

// Command quickstart is an example of using the Google Cloud Job Discovery API.
package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	talent "google.golang.org/api/jobs/v3"
)

func GetParent() string {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	return fmt.Sprintf("projects/%s", projectID)
}

func CreateCtsService() (*talent.Service, error) {
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}
	// Create the jobs service client.
	ctsService, err := talent.New(client)
	return ctsService, err
}

func main() {
	// Create the jobs service client.
	ctsService, err := CreateCtsService()
	if err != nil {
		log.Fatal(err)
	}

	// Make the RPC call.
	response, err := ctsService.Projects.Companies.List(DEFAULT_PROJECT_ID).Do()
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
