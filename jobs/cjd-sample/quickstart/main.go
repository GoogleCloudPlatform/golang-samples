// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START quick_start]
package main

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	jobs "google.golang.org/api/jobs/v2"
)

func main() {
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, jobs.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	// Create the jobs service client.
	jobService, err := jobs.New(client)
	if err != nil {
		log.Fatal(err)
	}

	// Make the RPC call.
	response, err := jobService.Companies.List().Do()
	if err != nil {
		log.Fatalf("Failed to list Companies: %v", err)
	}

	// Print the returned companies.
	for _, company := range response.Companies {
		fmt.Printf("================\nCompany:\n%q\n================\n", company)
	}
}

// [END quick_start]
