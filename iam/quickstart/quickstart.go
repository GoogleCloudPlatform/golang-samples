// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START iam_quickstart]

// The quickstart command is an example of using the IAM API.
package main

import (
	"log"

	"context"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/iam/v1"
)

func main() {
	// Get credentials.
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		log.Fatalf("google.DefaultClient: %v", err)
	}

	// Create the Cloud IAM service object.
	iamService, err := iam.New(client)
	if err != nil {
		log.Fatalf("iam.New: %v", err)
	}

	// Call the Cloud IAM Roles API.
	resp, err := iamService.Roles.List().Do()
	if err != nil {
		log.Fatalf("List.Do: %v", err)
	}

	// Process the response.
	for _, role := range resp.Roles {
		log.Println("Tile: " + role.Title)
		log.Println("Name: " + role.Name)
		log.Println("Description: " + role.Description)
		log.Println()
	}
}

// [END iam_quickstart]
