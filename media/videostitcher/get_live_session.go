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

// [START videostitcher_get_live_session]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	"cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
)

// getLiveSession gets a livestream session by ID.
func getLiveSession(w io.Writer, projectID, sessionID string) error {
	// projectID := "my-project-id"
	// sessionID := "123-456-789"
	location := "us-central1"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %w", err)
	}
	defer client.Close()

	req := &stitcherpb.GetLiveSessionRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/liveSessions/%s", projectID, location, sessionID),
	}
	// Gets the session.
	response, err := client.GetLiveSession(ctx, req)
	if err != nil {
		return fmt.Errorf("client.GetLiveSession: %w", err)
	}

	fmt.Fprintf(w, "Live session: %+v", response)
	return nil
}

// [END videostitcher_get_live_session]
