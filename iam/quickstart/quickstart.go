// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START iam_quickstart]

// The quickstart command is an example of using the IAM API.
package main

import (
	"fmt"

	"context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iam/v1"
)

func main() {
	// Get credentials
	client, _ := google.DefaultClient(
		context.Background(),
		iam.CloudPlatformScope)

	// Create the Cloud IAM service object
	iamService, err := iam.New(client)
	if err != nil {
		log.Fatalf("iam.New: %v", err)
	}

	// Call the Cloud IAM Roles API
	resp, err := iamService.Roles.List().Do()
	if err != nil {
		log.Fatalf("List.Do: %v", err)
	}
	roles := response.Roles

	// Process the response
	for _, role := range resp.Roles {
		fmt.Println("Tile: " + role.Title)
		fmt.Println("Name: " + role.Name)
		fmt.Println("Description: " + role.Description)
		fmt.Println()
	}
}

// [END iam_quickstart]
