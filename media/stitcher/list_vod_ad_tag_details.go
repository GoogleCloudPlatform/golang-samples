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

// [START video_stitcher_list_vod_ad_tag_details]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/api/iterator"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "google.golang.org/genproto/googleapis/cloud/video/stitcher/v1"
)

// listVodAdTagDetails lists the ad tag details for the specified VOD session.
func listVodAdTagDetails(w io.Writer, projectID, location, sessionID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// sessionID := "my-session-id"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	req := &stitcherstreampb.ListVodAdTagDetailsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/vodSessions/%s", projectID, location, sessionID),
	}

	it := client.ListVodAdTagDetails(ctx, req)
	fmt.Fprintln(w, "VOD ad tag details:")

	for {
		response, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("ListVodAdTagDetails: %v", err)
		}
		fmt.Fprintln(w, response.GetName())
	}
	return nil
}

// [END video_stitcher_list_vod_ad_tag_details]
