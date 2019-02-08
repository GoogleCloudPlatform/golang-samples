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

// [START quick_start]

// Command quickstart is an example of using the Google Cloud Talent Solution API.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	talent "google.golang.org/api/jobs/v3"
)

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	parent := fmt.Sprintf("projects/%s", projectID)

	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, talent.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	// Create the jobs service client.
	ctsService, err := talent.New(client)
	if err != nil {
		log.Fatal(err)
	}

	// Make the RPC call.
	response, err := ctsService.Projects.Companies.List(parent).Do()
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
