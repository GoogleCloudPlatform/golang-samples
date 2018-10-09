// Copyright 2018 Google LLC
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
        "fmt"
        "log"
        "os"

        "cloud.google.com/go/asset/v1beta1"
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
        gcsDestination := &assetpb.GcsDestination{
                Uri: string("gs://[my_gcs_bucket]/[my_asset_dump_file]"),
        }
        destination := &assetpb.OutputConfig_GcsDestination{
                GcsDestination: gcsDestination,
        }
        outputConfig := &assetpb.OutputConfig{
                Destination: destination,
        }
        req := &assetpb.ExportAssetsRequest{
                Parent:       fmt.Sprintf("projects/%s", projectID),
                OutputConfig: outputConfig,
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
