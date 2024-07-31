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
	stitcherstreampb "cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
)

// createVodSession creates a video on demand (VOD) session in which to insert ads.
// VOD sessions are ephemeral resources that expire after a few hours.
func createVodSession(w io.Writer, projectID, vodConfigID string) error {
	// projectID := "my-project-id"
	// vodConfigID := "my-vod-config-id"
	location := "us-central1"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %w", err)
	}
	defer client.Close()

	req := &stitcherstreampb.CreateVodSessionRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		VodSession: &stitcherstreampb.VodSession{
			VodConfig:  fmt.Sprintf("projects/%s/locations/%s/vodConfigs/%s", projectID, location, vodConfigID),
			AdTracking: stitcherstreampb.AdTracking_SERVER,
		},
	}
	// Creates the VOD session.
	response, err := client.CreateVodSession(ctx, req)
	if err != nil {
		return fmt.Errorf("client.CreateVodSession: %w", err)
	}

	fmt.Fprintf(w, "VOD session: %v", response.GetName())
	return nil
}

// [END videostitcher_create_vod_session]
