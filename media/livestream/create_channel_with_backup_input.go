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

package livestream

// [START livestream_create_channel_with_backup_input]
import (
	"context"
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes/duration"

	livestream "cloud.google.com/go/video/livestream/apiv1"
	"cloud.google.com/go/video/livestream/apiv1/livestreampb"
)

// createChannelWithBackupInput creates a channel with a failover backup input.
func createChannelWithBackupInput(w io.Writer, projectID, location, channelID, primaryInputID, backupInputID, outputURI string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// channelID := "my-channel"
	// primaryInputID := "my-primary-input"
	// backupInputID := "my-backup-input"
	// outputURI := "gs://my-bucket/my-output-folder/"
	ctx := context.Background()
	client, err := livestream.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	primaryInput := fmt.Sprintf("projects/%s/locations/%s/inputs/%s", projectID, location, primaryInputID)
	backupInput := fmt.Sprintf("projects/%s/locations/%s/inputs/%s", projectID, location, backupInputID)
	automaticFailover := &livestreampb.InputAttachment_AutomaticFailover{
		InputKeys: []string{"my-backup-input"},
	}

	req := &livestreampb.CreateChannelRequest{
		Parent:    fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		ChannelId: channelID,
		Channel: &livestreampb.Channel{
			InputAttachments: []*livestreampb.InputAttachment{
				{
					Key:               "my-primary-input",
					Input:             primaryInput,
					AutomaticFailover: automaticFailover,
				},
				{
					Key:   "my-backup-input",
					Input: backupInput,
				},
			},
			Output: &livestreampb.Channel_Output{
				Uri: outputURI,
			},
			ElementaryStreams: []*livestreampb.ElementaryStream{
				{
					Key: "es_video",
					ElementaryStream: &livestreampb.ElementaryStream_VideoStream{
						VideoStream: &livestreampb.VideoStream{
							CodecSettings: &livestreampb.VideoStream_H264{
								H264: &livestreampb.VideoStream_H264CodecSettings{
									Profile:      "high",
									BitrateBps:   3000000,
									FrameRate:    30,
									HeightPixels: 720,
									WidthPixels:  1280,
								},
							},
						},
					},
				},
				{
					Key: "es_audio",
					ElementaryStream: &livestreampb.ElementaryStream_AudioStream{
						AudioStream: &livestreampb.AudioStream{
							Codec:        "aac",
							ChannelCount: 2,
							BitrateBps:   160000,
						},
					},
				},
			},
			MuxStreams: []*livestreampb.MuxStream{
				{
					Key:               "mux_video",
					ElementaryStreams: []string{"es_video"},
					SegmentSettings: &livestreampb.SegmentSettings{
						SegmentDuration: &duration.Duration{
							Seconds: 2,
						},
					},
				},
				{
					Key:               "mux_audio",
					ElementaryStreams: []string{"es_audio"},
					SegmentSettings: &livestreampb.SegmentSettings{
						SegmentDuration: &duration.Duration{
							Seconds: 2,
						},
					},
				},
			},
			Manifests: []*livestreampb.Manifest{
				{
					FileName:        "manifest.m3u8",
					Type:            livestreampb.Manifest_HLS,
					MuxStreams:      []string{"mux_video", "mux_audio"},
					MaxSegmentCount: 5,
				},
			},
		},
	}
	// Creates the channel.
	op, err := client.CreateChannel(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateChannel: %w", err)
	}
	response, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %w", err)
	}

	fmt.Fprintf(w, "Channel: %v", response.Name)
	return nil
}

// [END livestream_create_channel_with_backup_input]
