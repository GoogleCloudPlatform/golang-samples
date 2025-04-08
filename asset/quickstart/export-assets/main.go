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

// [START asset_quickstart_export_assets]

// Sample asset-quickstart exports assets to given path.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
)

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	client, err := asset.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	bucketName := fmt.Sprintf("%s-for-assets", projectID)
	assetDumpFile := fmt.Sprintf("gs://%s/my-assets.txt", bucketName)
	req := &assetpb.ExportAssetsRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		OutputConfig: &assetpb.OutputConfig{
			Destination: &assetpb.OutputConfig_GcsDestination{
				GcsDestination: &assetpb.GcsDestination{
					ObjectUri: &assetpb.GcsDestination_Uri{
						Uri: string(assetDumpFile),
					},
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
