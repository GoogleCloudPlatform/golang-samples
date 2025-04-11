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

// [START videostitcher_create_vod_config]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
)

// createVodConfig creates a VOD config. VOD configs are used to configure VOD sessions.
func createVodConfig(w io.Writer, projectID, vodConfigID, sourceURI string) error {
	// projectID := "my-project-id"
	// vodConfigID := "my-vod-config-id"

	// Uri of the media to stitch; this URI must reference either an MPEG-DASH
	// manifest (.mpd) file or an M3U playlist manifest (.m3u8) file.
	// sourceURI := "https://storage.googleapis.com/my-bucket/main.mpd"

	// See https://cloud.google.com/video-stitcher/docs/concepts for information
	// on ad tags and ad metadata. This sample uses an ad tag URL that displays
	// a VMAP Pre-roll ad
	// (https://developers.google.com/interactive-media-ads/docs/sdks/html5/client-side/tags).
	adTagURI := "https://pubads.g.doubleclick.net/gampad/ads?iu=/21775744923/external/vmap_ad_samples&sz=640x480&cust_params=sample_ar%3Dpreonly&ciu_szs=300x250%2C728x90&gdfp_req=1&ad_rule=1&output=vmap&unviewed_position_start=1&env=vp&impl=s&correlator="
	location := "us-central1"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %w", err)
	}
	defer client.Close()

	req := &stitcherstreampb.CreateVodConfigRequest{
		Parent:      fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		VodConfigId: vodConfigID,
		VodConfig: &stitcherstreampb.VodConfig{
			SourceUri: sourceURI,
			AdTagUri:  adTagURI,
		},
	}
	// Creates the VOD config.
	op, err := client.CreateVodConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("client.CreateVodConfig: %w", err)
	}
	response, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "VOD config: %v", response.GetName())
	return nil
}

// [END videostitcher_create_vod_config]
