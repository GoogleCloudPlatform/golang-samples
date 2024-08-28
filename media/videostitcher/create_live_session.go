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

// [START videostitcher_create_live_session]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	"cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
)

// createLiveSession creates a livestream session in which to insert ads.
// Live sessions are ephemeral resources that expire after a few minutes.
func createLiveSession(w io.Writer, projectID, liveConfigID string) error {
	// projectID := "my-project-id"
	// liveConfigID := "my-live-config"
	location := "us-central1"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %w", err)
	}
	defer client.Close()

	req := &stitcherpb.CreateLiveSessionRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		LiveSession: &stitcherpb.LiveSession{
			LiveConfig: fmt.Sprintf("projects/%s/locations/%s/liveConfigs/%s", projectID, location, liveConfigID),
		},
	}
	// Creates the live session.
	response, err := client.CreateLiveSession(ctx, req)
	if err != nil {
		return fmt.Errorf("client.CreateLiveSession: %w", err)
	}

	fmt.Fprintf(w, "Live session: %v\n", response.GetName())
	fmt.Fprintf(w, "Play URI: %v", response.GetPlayUri())
	return nil
}

// [END videostitcher_create_live_session]
