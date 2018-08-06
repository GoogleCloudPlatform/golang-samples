// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START quick_start]

// Command quickstart is an example of using the Google Cloud Job Discovery API.
package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	jobs "google.golang.org/api/jobs/v3"
)

var DEFAULT_PROJECT_ID = "projects/" + os.Getenv("GOOGLE_CLOUD_PROJECT")

func main() {
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, jobs.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	// Create the jobs service client.
	service, err := jobs.New(client)
	if err != nil {
		log.Fatal(err)
	}
	if service.Projects == nil {
		log.Fatal(errors.New("Failed to initiate Projects"))
	}
	if service.Projects.Companies == nil {
		log.Fatal(errors.New("Failed to initiate Companies service"))
	}

	// Make the RPC call.
	response, err := service.Projects.Companies.List(DEFAULT_PROJECT_ID).Do()
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
