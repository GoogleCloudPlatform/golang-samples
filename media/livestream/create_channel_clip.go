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

package livestream

// [START livestream_create_channel_clip]
import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"

	livestream "cloud.google.com/go/video/livestream/apiv1"
	"cloud.google.com/go/video/livestream/apiv1/livestreampb"
)

// createChannelClip creates a channel clip. A clip is a sub-resource of a
// channel. You can use a channel clip to create video on demand (VOD) files
// from a live stream. These VOD files are saved to Cloud Storage.
func createChannelClip(w io.Writer, projectID, channelID, clipID, outputURI string) error {
	// projectID := "my-project-id"
	// channelID := "my-channel"
	// clipID := "my-channel-clip"
	// outputURI := "gs://my-bucket/my-output-folder/"
	location := "us-central1"
	ctx := context.Background()
	client, err := livestream.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &livestreampb.CreateClipRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/channels/%s", projectID, location, channelID),
		ClipId: clipID,
		Clip: &livestreampb.Clip{
			OutputUri: outputURI,
			Slices: []*livestreampb.Clip_Slice{
				{
					Kind: &livestreampb.Clip_Slice_TimeSlice{
						TimeSlice: &livestreampb.Clip_TimeSlice{
							// Create a 20 second clip starting 40 seconds ago
							MarkinTime: &timestamp.Timestamp{
								Seconds: time.Now().Unix() - 40,
							},
							MarkoutTime: &timestamp.Timestamp{
								Seconds: time.Now().Unix() - 20,
							},
						},
					},
				},
			},
			ClipManifests: []*livestreampb.Clip_ClipManifest{
				{
					ManifestKey: "manifest_hls",
				},
			},
		},
	}
	// Creates the channel clip.
	op, err := client.CreateClip(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateClip: %w", err)
	}
	response, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Channel clip: %v", response.GetName())
	return nil
}

// [END livestream_create_channel_clip]
