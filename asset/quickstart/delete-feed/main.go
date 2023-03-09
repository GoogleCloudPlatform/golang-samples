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

// [START asset_quickstart_delete_feed]

// Sample delete-feed delete feed.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

// Command-line flags.
var (
	feedID = flag.String("feed_id", "YOUR_FEED_ID", "Identifier of Feed.")
)

func main() {
	flag.Parse()
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		log.Fatalf("asset.NewClient: %v", err)
	}
	defer client.Close()

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
	feedName := fmt.Sprintf("projects/%s/feeds/%s", projectNumber, *feedID)
	req := &assetpb.DeleteFeedRequest{
		Name: feedName,
	}
	if err = client.DeleteFeed(ctx, req); err != nil {
		log.Fatalf("client.DeleteFeed: %v", err)
	}
	fmt.Print("Deleted Feed")
}

// [END asset_quickstart_delete_feed]
