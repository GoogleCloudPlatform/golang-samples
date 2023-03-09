// Copyright 2020 Google LLC
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

// [START asset_quickstart_list_assets]

// Sample list-assets list assets.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/iterator"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
)

func main() {
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	assetType := "storage.googleapis.com/Bucket"
	req := &assetpb.ListAssetsRequest{
		Parent:      fmt.Sprintf("projects/%s", projectID),
		AssetTypes:  []string{assetType},
		ContentType: assetpb.ContentType_RESOURCE,
	}

	// Call ListAssets API to get an asset iterator.
	it := client.ListAssets(ctx, req)

	// Traverse and print the first 10 listed assets in response.
	for i := 0; i < 10; i++ {
		response, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(response)
	}
}

// [END asset_quickstart_list_assets]
