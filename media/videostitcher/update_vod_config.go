// Copyright 2024 Google LLC
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

package videostitcher

// [START videostitcher_update_vod_config]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateVodConfig updates an existing VOD config. This sample updates the sourceURI for an
// existing VOD config.
func updateVodConfig(w io.Writer, projectID, vodConfigID, sourceURI string) error {
	// projectID := "my-project-id"
	// vodConfigID := "my-vod-config-id"
	// sourceURI := "https://storage.googleapis.com/my-bucket/main.mpd"
	location := "us-central1"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %w", err)
	}
	defer client.Close()

	req := &stitcherstreampb.UpdateVodConfigRequest{
		VodConfig: &stitcherstreampb.VodConfig{
			Name:      fmt.Sprintf("projects/%s/locations/%s/vodConfigs/%s", projectID, location, vodConfigID),
			SourceUri: sourceURI,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"sourceUri",
			},
		},
	}
	// Updates the VOD config.
	op, err := client.UpdateVodConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("client.UpdateVodConfig: %w", err)
	}
	response, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Updated VOD config: %+v", response)
	return nil
}

// [END videostitcher_update_vod_config]
