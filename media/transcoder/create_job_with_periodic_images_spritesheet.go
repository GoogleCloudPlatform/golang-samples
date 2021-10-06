// Copyright 2021 Google LLC
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

package transcoder

// [START transcoder_create_job_with_periodic_images_spritesheet]
import (
	"context"
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes/duration"

	transcoder "cloud.google.com/go/video/transcoder/apiv1"
	transcoderpb "google.golang.org/genproto/googleapis/cloud/video/transcoder/v1"
)

// createJobWithPeriodicImagesSpritesheet creates a job from an ad-hoc configuration and generates
// two spritesheets from the input video. Each spritesheet contains images that are captured
// periodically based on a user-defined time interval.
func createJobWithPeriodicImagesSpritesheet(w io.Writer, projectID string, location string, inputURI string, outputURI string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// inputURI := "gs://my-bucket/my-video-file"
	// outputURI := "gs://my-bucket/my-output-folder/"
	ctx := context.Background()
	client, err := transcoder.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	req := &transcoderpb.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Job: &transcoderpb.Job{
			InputUri:  inputURI,
			OutputUri: outputURI,
			JobConfig: &transcoderpb.Job_Config{
				Config: &transcoderpb.JobConfig{
					ElementaryStreams: []*transcoderpb.ElementaryStream{
						&transcoderpb.ElementaryStream{
							Key: "video_stream0",
							ElementaryStream: &transcoderpb.ElementaryStream_VideoStream{
								VideoStream: &transcoderpb.VideoStream{
									CodecSettings: &transcoderpb.VideoStream_H264{
										H264: &transcoderpb.VideoStream_H264CodecSettings{
											BitrateBps:   550000,
											FrameRate:    60,
											HeightPixels: 360,
											WidthPixels:  640,
										},
									},
								},
							},
						},
						&transcoderpb.ElementaryStream{
							Key: "audio_stream0",
							ElementaryStream: &transcoderpb.ElementaryStream_AudioStream{
								AudioStream: &transcoderpb.AudioStream{
									Codec:      "aac",
									BitrateBps: 64000,
								},
							},
						},
					},
					MuxStreams: []*transcoderpb.MuxStream{
						&transcoderpb.MuxStream{
							Key:               "sd",
							Container:         "mp4",
							ElementaryStreams: []string{"video_stream0", "audio_stream0"},
						},
					},
					SpriteSheets: []*transcoderpb.SpriteSheet{
						&transcoderpb.SpriteSheet{
							FilePrefix:         "small-sprite-sheet",
							SpriteWidthPixels:  64,
							SpriteHeightPixels: 32,
							ExtractionStrategy: &transcoderpb.SpriteSheet_Interval{
								Interval: &duration.Duration{
									Seconds: 7,
								},
							},
						},
						&transcoderpb.SpriteSheet{
							FilePrefix:         "large-sprite-sheet",
							SpriteWidthPixels:  128,
							SpriteHeightPixels: 72,
							ExtractionStrategy: &transcoderpb.SpriteSheet_Interval{
								Interval: &duration.Duration{
									Seconds: 7,
								},
							},
						},
					},
				},
			},
		},
	}
	// Creates the job. Jobs take a variable amount of time to run. You can query for the job state.
	// See https://cloud.google.com/transcoder/docs/how-to/jobs#check_job_status for more info.
	response, err := client.CreateJob(ctx, req)
	if err != nil {
		return fmt.Errorf("createJobWithPeriodicImagesSpritesheet: %v", err)
	}

	fmt.Fprintf(w, "Job: %v", response.GetName())
	return nil
}

// [END transcoder_create_job_with_periodic_images_spritesheet]
