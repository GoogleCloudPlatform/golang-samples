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

package stitcher

// [START video_stitcher_get_vod_ad_tag_detail]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "google.golang.org/genproto/googleapis/cloud/video/stitcher/v1"
)

// getVodAdTagDetail gets the specified ad tag detail for a VOD session.
func getVodAdTagDetail(w io.Writer, projectID, location, sessionID, adTagDetailID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// sessionID := "my-session-id"
	// adTagDetailID := "my-ad-tag-detail-id"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	req := &stitcherstreampb.GetVodAdTagDetailRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/vodSessions/%s/vodAdTagDetails/%s", projectID, location, sessionID, adTagDetailID),
	}

	response, err := client.GetVodAdTagDetail(ctx, req)
	if err != nil {
		return fmt.Errorf("GetVodAdTagDetail: %v", err)
	}

	fmt.Fprintf(w, "VOD ad tag detail: %v", response.Name)
	return nil
}

// [END video_stitcher_get_vod_ad_tag_detail]
