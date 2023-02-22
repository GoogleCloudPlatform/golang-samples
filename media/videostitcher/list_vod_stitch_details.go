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

// [START videostitcher_list_vod_stitch_details]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/api/iterator"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	"cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
)

// listVodStitchDetails lists the stitch details for a video on demand (VOD) session.
func listVodStitchDetails(w io.Writer, projectID, sessionID string) error {
	// projectID := "my-project-id"
	// sessionID := "123-456-789"
	location := "us-central1"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	req := &stitcherpb.ListVodStitchDetailsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/vodSessions/%s", projectID, location, sessionID),
	}

	it := client.ListVodStitchDetails(ctx, req)
	fmt.Fprintln(w, "VOD stitch details:")
	for {
		response, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("it.Next(): %v", err)
		}
		fmt.Fprintln(w, response.GetName())
	}

	return nil
}

// [END videostitcher_list_vod_stitch_details]
