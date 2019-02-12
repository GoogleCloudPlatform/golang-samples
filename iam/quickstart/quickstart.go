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

// [START iam_quickstart]

// The quickstart command is an example of using the Cloud IAM Roles API.
package main

import (
	"context"
	"log"

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
	service, err := iam.New(client)
	if err != nil {
		log.Fatalf("iam.New: %v", err)
	}

	// Call the Cloud IAM Roles API.
	resp, err := service.Roles.List().Do()
	if err != nil {
		log.Fatalf("Roles.List: %v", err)
	}

	// Process the response.
	for _, role := range resp.Roles {
		log.Println("Tile: " + role.Title)
		log.Println("Name: " + role.Name)
		log.Println("Description: " + role.Description)
	}
}

// [END iam_quickstart]
