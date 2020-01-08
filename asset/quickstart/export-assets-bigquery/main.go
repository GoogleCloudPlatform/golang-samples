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
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
)

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	client, err := asset.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	datasetID := strings.ReplaceAll(fmt.Sprintf("%s-for-assets", projectID), "-", "_")
	dataset := fmt.Sprintf("projects/%s/datasets/%s", projectID, datasetID)
	table := "test"
	req := &assetpb.ExportAssetsRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		OutputConfig: &assetpb.OutputConfig{
			Destination: &assetpb.OutputConfig_BigqueryDestination{
				BigqueryDestination: &assetpb.BigQueryDestination{
					Dataset: string(dataset),
					Table:   string(table),
				},
			},
		},
	}
	operation_bq, err := client.ExportAssets(ctx, req_bq)
	if err != nil {
		log.Fatal(err)
	}
	response_bq, err := operation_bq.Wait(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(response_bq)
}

// [END asset_quickstart_export_assets_bigquery]
