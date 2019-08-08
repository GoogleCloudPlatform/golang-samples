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

// [START asset_quickstart_get_feed]

// Sample asset-quickstart get feed.

package main

import (
        "context"
        "fmt"
        "log"
        "os"
        
        assetUtils "github.com/GoogleCloudPlatform/golang-samples/asset/utils"
        asset "cloud.google.com/go/asset/apiv1p2beta1"
        assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1p2beta1"
)

func main() {
        ctx := context.Background()
        client, err := asset.NewClient(ctx)
        if err != nil {
                log.Fatal(err)
        }

        projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
        projectNumber := assetUtils.GetProjectNumberById(projectID)
        feedName := fmt.Sprintf("projects/%s/feeds/%s", projectNumber, "YOUR_FEED_ID")
        req := &assetpb.GetFeedRequest{
                Name: feedName}
        response, err := client.GetFeed(ctx, req)
        if err != nil {
                log.Fatal(err)
        }
        fmt.Print(response)
}
// [END asset_quickstart_get_feed]