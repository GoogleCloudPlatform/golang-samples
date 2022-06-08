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

// [START video_stitcher_get_vod_ad_tag_detail]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherpb "google.golang.org/genproto/googleapis/cloud/video/stitcher/v1"
)

// getVodAdTagDetail gets the specified ad tag detail for a video on demand (VOD) session.
func getVodAdTagDetail(w io.Writer, projectID, sessionID, adTagDetailID string) error {
	// projectID := "my-project-id"
	// sessionID := "123-456-789"
	// adTagDetailID := "01234-56789"
	location := "us-central1"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	req := &stitcherpb.GetVodAdTagDetailRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/vodSessions/%s/vodAdTagDetails/%s", projectID, location, sessionID, adTagDetailID),
	}
	// Gets the ad tag detail.
	response, err := client.GetVodAdTagDetail(ctx, req)
	if err != nil {
		return fmt.Errorf("client.GetVodAdTagDetail: %v", err)
	}
	b, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %v", err)
	}

	fmt.Fprintf(w, "VOD ad tag detail:\n%s", string(b))
	return nil
}

// [END video_stitcher_get_vod_ad_tag_detail]
