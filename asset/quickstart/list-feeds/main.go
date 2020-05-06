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

// [START asset_quickstart_list_feeds]

// Sample list-feeds list feeds.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	asset "cloud.google.com/go/asset/apiv1"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
)

func main() {
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		log.Fatalf("asset.NewClient: %v", err)
	}

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	cloudresourcemanagerClient, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		log.Fatalf("cloudresourcemanager.NewService: %v", err)
	}

	project, err := cloudresourcemanagerClient.Projects.Get(projectID).Do()
	if err != nil {
		log.Fatalf("cloudresourcemanagerClient.Projects.Get.Do: %v", err)
	}
	projectNumber := strconv.FormatInt(project.ProjectNumber, 10)
	parent := fmt.Sprintf("projects/%s", projectNumber)
	req := &assetpb.ListFeedsRequest{
		Parent: parent}
	response, err := client.ListFeeds(ctx, req)
	if err != nil {
		log.Fatalf("client.ListFeeds: %v", err)
	}
	fmt.Print(response)
}

// [END asset_quickstart_list_feeds]
