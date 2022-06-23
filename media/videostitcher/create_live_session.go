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

// [START video_stitcher_create_live_session]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherpb "google.golang.org/genproto/googleapis/cloud/video/stitcher/v1"
)

// createLiveSession creates a livestream session in which to insert ads.
// Live sessions are ephemeral resources that expire after a few minutes.
func createLiveSession(w io.Writer, projectID, sourceURI, slateID string) error {
	// projectID := "my-project-id"

	// Uri of the media to stitch; this URI must reference either an MPEG-DASH
	// manifest (.mpd) file or an M3U playlist manifest (.m3u8) file.
	// sourceURI := "https://storage.googleapis.com/my-bucket/main.mpd"
	// slateID := "my-slate"
	// See https://cloud.google.com/video-stitcher/docs/concepts for information
	// on ad tags and ad metadata. This sample uses an ad tag URL that displays
	// a Single Inline Linear ad
	// (https://developers.google.com/interactive-media-ads/docs/sdks/html5/client-side/tags).
	adTagURI := "https://pubads.g.doubleclick.net/gampad/ads?iu=/21775744923/external/single_ad_samples&sz=640x480&cust_params=sample_ct%3Dlinear&ciu_szs=300x250%2C728x90&gdfp_req=1&output=vast&unviewed_position_start=1&env=vp&impl=s&correlator="
	location := "us-central1"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	req := &stitcherpb.CreateLiveSessionRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		LiveSession: &stitcherpb.LiveSession{
			SourceUri:      sourceURI,
			DefaultAdTagId: "default",
			DefaultSlateId: slateID,
			AdTagMap: map[string]*stitcherpb.AdTag{"default": &stitcherpb.AdTag{
				Uri: adTagURI,
			}},
		},
	}
	// Creates the live session.
	response, err := client.CreateLiveSession(ctx, req)
	if err != nil {
		return fmt.Errorf("client.CreateLiveSession: %v", err)
	}

	fmt.Fprintf(w, "Live session: %v\n", response.GetName())
	fmt.Fprintf(w, "Play URI: %v", response.GetPlayUri())
	return nil
}

// [END video_stitcher_create_live_session]
