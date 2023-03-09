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

// [START asset_quickstart_batch_get_assets_history]

// Sample asset-quickstart batch-gets assets history.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func main() {
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	bucketResourceName := fmt.Sprintf("//storage.googleapis.com/%s-for-assets", projectID)
	req := &assetpb.BatchGetAssetsHistoryRequest{
		Parent:      fmt.Sprintf("projects/%s", projectID),
		AssetNames:  []string{bucketResourceName},
		ContentType: assetpb.ContentType_RESOURCE,
		ReadTimeWindow: &assetpb.TimeWindow{
			StartTime: &timestamp.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
	}
	response, err := client.BatchGetAssetsHistory(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(response)
}

// [END asset_quickstart_batch_get_assets_history]
