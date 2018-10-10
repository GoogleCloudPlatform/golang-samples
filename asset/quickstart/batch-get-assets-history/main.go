// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START asset_quickstart_batch_get_assets_history]

// Sample asset-quickstart batch-gets assets history.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	asset "cloud.google.com/go/asset/v1beta1"
	"github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/net/context"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1beta1"
)

func main() {
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

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
