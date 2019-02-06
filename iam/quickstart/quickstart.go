// Copyright 2019 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
