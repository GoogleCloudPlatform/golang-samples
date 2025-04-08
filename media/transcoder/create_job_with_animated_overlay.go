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

// [START transcoder_create_job_with_animated_overlay]
import (
	"context"
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes/duration"

	transcoder "cloud.google.com/go/video/transcoder/apiv1"
	"cloud.google.com/go/video/transcoder/apiv1/transcoderpb"
)

// createJobWithAnimatedOverlay creates a job based on a given configuration that
// includes an animated overlay. See
// https://cloud.google.com/transcoder/docs/how-to/create-overlays#create-animated-overlay
// for more information.
func createJobWithAnimatedOverlay(w io.Writer, projectID string, location string, inputURI string, overlayImageURI string, outputURI string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// inputURI := "gs://my-bucket/my-video-file"
	// overlayImageURI := "gs://my-bucket/my-overlay-image-file"
	// outputURI := "gs://my-bucket/my-output-folder/"
	ctx := context.Background()
	client, err := transcoder.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
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
						{
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
						{
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
						{
							Key:               "sd",
							Container:         "mp4",
							ElementaryStreams: []string{"video_stream0", "audio_stream0"},
						},
					},
					Overlays: []*transcoderpb.Overlay{
						{
							Image: &transcoderpb.Overlay_Image{
								Uri: overlayImageURI,
								Resolution: &transcoderpb.Overlay_NormalizedCoordinate{
									X: 0,
									Y: 0,
								},
								Alpha: 1,
							},
							Animations: []*transcoderpb.Overlay_Animation{
								{
									AnimationType: &transcoderpb.Overlay_Animation_AnimationFade{
										AnimationFade: &transcoderpb.Overlay_AnimationFade{
											FadeType: transcoderpb.Overlay_FADE_IN,
											Xy: &transcoderpb.Overlay_NormalizedCoordinate{
												X: 0.5,
												Y: 0.5,
											},
											StartTimeOffset: &duration.Duration{
												Seconds: 5,
											},
											EndTimeOffset: &duration.Duration{
												Seconds: 10,
											},
										},
									},
								},

								{
									AnimationType: &transcoderpb.Overlay_Animation_AnimationFade{
										AnimationFade: &transcoderpb.Overlay_AnimationFade{
											FadeType: transcoderpb.Overlay_FADE_OUT,
											Xy: &transcoderpb.Overlay_NormalizedCoordinate{
												X: 0.5,
												Y: 0.5,
											},
											StartTimeOffset: &duration.Duration{
												Seconds: 12,
											},
											EndTimeOffset: &duration.Duration{
												Seconds: 15,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	// Creates the job. Jobs take a variable amount of time to run.
	// You can query for the job state; see getJob() in get_job.go.
	response, err := client.CreateJob(ctx, req)
	if err != nil {
		return fmt.Errorf("createJobWithAnimatedOverlay: %w", err)
	}

	fmt.Fprintf(w, "Job: %v", response.GetName())
	return nil
}

// [END transcoder_create_job_with_animated_overlay]
