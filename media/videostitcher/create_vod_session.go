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

package videostitcher

// [START videostitcher_create_vod_session]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherpb "google.golang.org/genproto/googleapis/cloud/video/stitcher/v1"
)

// createVodSession creates a video on demand (VOD) session in which to insert ads.
// VOD sessions are ephemeral resources that expire after a few hours.
func createVodSession(w io.Writer, projectID, sourceURI string) error {
	// projectID := "my-project-id"

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
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	req := &stitcherpb.CreateVodSessionRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		VodSession: &stitcherpb.VodSession{
			SourceUri: sourceURI,
			AdTagUri:  adTagURI,
		},
	}
	// Creates the VOD session.
	response, err := client.CreateVodSession(ctx, req)
	if err != nil {
		return fmt.Errorf("client.CreateVodSession: %v", err)
	}

	fmt.Fprintf(w, "VOD session: %v", response.GetName())
	return nil
}

// [END videostitcher_create_vod_session]
