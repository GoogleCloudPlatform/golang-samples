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

// [START videostitcher_get_vod_config]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
)

// getVodConfig gets a previously-created VOD config.
func getVodConfig(w io.Writer, projectID, vodConfigID string) error {
	// projectID := "my-project-id"
	// vodConfigID := "my-vod-config-id"
	location := "us-central1"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %w", err)
	}
	defer client.Close()

	req := &stitcherstreampb.GetVodConfigRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/vodConfigs/%s", projectID, location, vodConfigID),
	}

	response, err := client.GetVodConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("client.GetVodConfig: %w", err)
	}
	b, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}

	fmt.Fprintf(w, "VOD config:\n%s", string(b))
	return nil
}

// [END videostitcher_get_vod_config]
