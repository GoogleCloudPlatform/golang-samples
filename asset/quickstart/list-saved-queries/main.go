// Copyright 2022 Google LLC
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

// [START asset_quickstart_list_saved_queries]

// Sample list-saved-queries list saved queries.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"google.golang.org/api/iterator"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

func main() {
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		log.Fatalf("asset.NewClient: %v", err)
	}
	defer client.Close()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	fmt.Println("projectID:", projectID)
	cloudresourcemanagerClient, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		log.Fatalf("cloudresourcemanager.NewService: %v", err)
	}

	project, err := cloudresourcemanagerClient.Projects.Get(projectID).Do()
	if err != nil {
		log.Fatalf("cloudresourcemanagerClient.Projects.Get.Do: %v", err)
	}
	projectNumber := strconv.FormatInt(project.ProjectNumber, 10)
	fmt.Println("projectNumber: %s", projectNumber)
	parent := fmt.Sprintf("projects/%s", projectNumber)
	req := &assetpb.ListSavedQueriesRequest{
		Parent: parent}
	it := client.ListSavedQueries(ctx, req)
	for {
		response, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(response)
	}
}

// [END asset_quickstart_list_saved_queries]