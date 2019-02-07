// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START iam_quickstart]

package quickstart

import (
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iam/v1"
)

func main() {
	// Get credentials
	client, _ := google.DefaultClient(
		context.Background(),
		iam.CloudPlatformScope)

	// Create the Cloud IAM service object
	iamService, _ := iam.New(client)

	// Call the Cloud IAM Roles API
	response, _ := iamService.Roles.List().Do()
	roles := response.Roles

	// Process the response
	for _, role := range roles {
		fmt.Println("Tile: " + role.Title)
		fmt.Println("Name: " + role.Name)
		fmt.Println("Description: " + role.Description)
		fmt.Println()
	}
}

// [END iam_quickstart]
