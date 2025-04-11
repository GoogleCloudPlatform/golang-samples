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

// [START videostitcher_get_vod_session]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
)

// getVodSession gets a VOD session by ID.
func getVodSession(w io.Writer, projectID, sessionID string) error {
	// projectID := "my-project-id"
	// sessionID := "123-456-789"
	location := "us-central1"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %w", err)
	}
	defer client.Close()

	req := &stitcherstreampb.GetVodSessionRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/vodSessions/%s", projectID, location, sessionID),
	}
	// Gets the session.
	response, err := client.GetVodSession(ctx, req)
	if err != nil {
		return fmt.Errorf("client.GetVodSession: %w", err)
	}
	b, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}

	fmt.Fprintf(w, "VOD session:\n%s", string(b))
	return nil
}

// [END videostitcher_get_vod_session]
