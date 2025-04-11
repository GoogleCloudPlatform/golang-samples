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

// [START asset_quickstart_create_feed]

// Sample create-feed create feed.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
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
	feedParent := fmt.Sprintf("projects/%s", projectID)
	assetNames := []string{"YOUR_ASSET_NAME"}
	topic := fmt.Sprintf("projects/%s/topics/%s", projectID, "YOUR_TOPIC_NAME")

	req := &assetpb.CreateFeedRequest{
		Parent: feedParent,
		FeedId: *feedID,
		Feed: &assetpb.Feed{
			AssetNames: assetNames,
			FeedOutputConfig: &assetpb.FeedOutputConfig{
				Destination: &assetpb.FeedOutputConfig_PubsubDestination{
					PubsubDestination: &assetpb.PubsubDestination{
						Topic: topic,
					},
				},
			},
		}}
	response, err := client.CreateFeed(ctx, req)
	if err != nil {
		log.Fatalf("client.CreateFeed: %v", err)
	}
	fmt.Print(response)
}

// [END asset_quickstart_create_feed]
