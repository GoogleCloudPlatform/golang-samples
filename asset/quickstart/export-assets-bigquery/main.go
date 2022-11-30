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

// [START asset_quickstart_export_assets_bigquery]

// Sample asset-quickstart exports assets to given bigquery table.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
)

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	client, err := asset.NewClient(ctx)
	if err != nil {
		log.Fatalf("asset.NewClient: %v", err)
	}
	defer client.Close()
	datasetID := strings.Replace(fmt.Sprintf("%s-for-assets", projectID), "-", "_", -1)
	dataset := fmt.Sprintf("projects/%s/datasets/%s", projectID, datasetID)
	req := &assetpb.ExportAssetsRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		OutputConfig: &assetpb.OutputConfig{
			Destination: &assetpb.OutputConfig_BigqueryDestination{
				BigqueryDestination: &assetpb.BigQueryDestination{
					Dataset: dataset,
					Table:   "test",
					Force:   true,
				},
			},
		},
	}
	op, err := client.ExportAssets(ctx, req)
	if err != nil {
		log.Fatalf("ExportAssets: %v", err)
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		log.Fatalf("Wait: %v", err)
	}
	fmt.Print(resp)
}

// [END asset_quickstart_export_assets_bigquery]
