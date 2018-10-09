// Copyright 2018 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START asset-quickstart_export_assets]
package main

import (
	"cloud.google.com/go/asset/v1beta1"
	"fmt"
	"golang.org/x/net/context"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1beta1"
)

func main() {
	context := context.Background()
	projectId := "[PROJECT]"
	client, error := asset.NewClient(context)
	if error != nil {
		// TODO: Handle error.
		fmt.Print(error)
		return
	}
	request := &assetpb.ExportAssetsRequest{
		Parent: fmt.Sprintf("projects/%s", projectId),
		OutputConfig: &assetpb.OutputConfig{
			Destination: &assetpb.OutputConfig_GcsDestination{
				GcsDestination: &assetpb.GcsDestination{
					Uri: string("gs://[my_gcs_bucket]/[my_asset_dump_file]"),
				},
			},
		},
	}
	operation, error := client.ExportAssets(context, request)
	if error != nil {
		// TODO: Handle error.
		fmt.Print(error)
		return
	}
	response, error := operation.Wait(context)
	if error != nil {
		// TODO: Handle error.
		fmt.Print(error)
		return
	}
	// Do things with response.
	fmt.Print(response)
}

// [END asset-quickstart_export_assets]
