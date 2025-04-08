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

// [START asset_quickstart_update_feed]

// Sample update-feed update feed.
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
	field_mask "google.golang.org/genproto/protobuf/field_mask"
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
	topic := fmt.Sprintf("projects/%s/topics/%s", projectID, "TOPIC_TO_UPDATE")

	req := &assetpb.UpdateFeedRequest{
		Feed: &assetpb.Feed{
			Name: feedName,
			FeedOutputConfig: &assetpb.FeedOutputConfig{
				Destination: &assetpb.FeedOutputConfig_PubsubDestination{
					PubsubDestination: &assetpb.PubsubDestination{
						Topic: topic,
					},
				},
			},
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"feed_output_config.pubsub_destination.topic"},
		},
	}
	response, err := client.UpdateFeed(ctx, req)
	if err != nil {
		log.Fatalf("client.UpdateFeed: %v", err)
	}
	fmt.Print(response)
}

// [END asset_quickstart_update_feed]
