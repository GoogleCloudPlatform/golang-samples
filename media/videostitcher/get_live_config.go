// Copyright 2023 Google LLC
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

// [START videostitcher_get_live_config]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
)

// getLiveConfig gets a previously-created live config.
func getLiveConfig(w io.Writer, projectID, liveConfigID string) error {
	// projectID := "my-project-id"
	// liveConfigID := "my-live-config-id"
	location := "us-central1"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %w", err)
	}
	defer client.Close()

	req := &stitcherstreampb.GetLiveConfigRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/liveConfigs/%s", projectID, location, liveConfigID),
	}

	response, err := client.GetLiveConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("client.GetLiveConfig: %w", err)
	}
	b, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}

	fmt.Fprintf(w, "Live config:\n%s", string(b))
	return nil
}

// [END videostitcher_get_live_config]
