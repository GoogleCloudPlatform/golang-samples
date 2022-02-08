// Copyright 2022 Google LLC
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

package main

// [START asset_quickstart_export_assets_bigquery]
import (
	"context"
	"fmt"

	asset "cloud.google.com/go/asset/apiv1"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
)

// exportAssetsToBigQuery demonstrates how to export assets to a given bigquery
// table.
func exportAssetsToBigQuery(projectID, datasetID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("asset.NewClient: %v", err)
	}
	defer client.Close()
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
		return fmt.Errorf("client.ExportAssets: %v", err)
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("op.Wait: %v", err)
	}
	fmt.Print(resp)
	return nil
}

// [END asset_quickstart_export_assets_bigquery]
