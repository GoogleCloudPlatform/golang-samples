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

// [START video_stitcher_get_cdn_key]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherpb "google.golang.org/genproto/googleapis/cloud/video/stitcher/v1"
)

// getCdnKey gets a CDN key by ID.
func getCdnKey(w io.Writer, projectID, cdnKeyID string) error {
	// projectID := "my-project-id"
	// cdnKeyID := "my-cdn-key"
	location := "us-central1"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	req := &stitcherpb.GetCdnKeyRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", projectID, location, cdnKeyID),
	}
	// Gets the CDN key.
	response, err := client.GetCdnKey(ctx, req)
	if err != nil {
		return fmt.Errorf("client.GetCdnKey: %v", err)
	}

	fmt.Fprintf(w, "CDN key: %+v", response)
	return nil
}

// [END video_stitcher_get_cdn_key]
