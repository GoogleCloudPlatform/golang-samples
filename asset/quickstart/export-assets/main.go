// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START asset_quickstart_export_assets]

// Sample asset-quickstart exports assets to given path.
package main

import (
	"fmt"
	"log"
	"os"

	asset "cloud.google.com/go/asset/v1beta1"
	"golang.org/x/net/context"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1beta1"
)

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	client, err := asset.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	bucketName := fmt.Sprintf("%s-for-assets", projectID)
	assetDumpFile := fmt.Sprintf("gs://%s/my-assets.txt", bucketName)
	req := &assetpb.ExportAssetsRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		OutputConfig: &assetpb.OutputConfig{
			Destination: &assetpb.OutputConfig_GcsDestination{
				GcsDestination: &assetpb.GcsDestination{
					Uri: string(assetDumpFile),
				},
			},
		},
	}
	operation, err := client.ExportAssets(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	response, err := operation.Wait(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(response)
}

// [END asset_quickstart_export_assets]
