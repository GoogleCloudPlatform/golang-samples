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

// [START transcoder_create_job_with_static_overlay]
import (
	"context"
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes/duration"

	transcoder "cloud.google.com/go/video/transcoder/apiv1beta1"
	transcoderpb "google.golang.org/genproto/googleapis/cloud/video/transcoder/v1beta1"
)

// createJobWithStaticOverlay creates a job based on a given configuration that
// includes a static overlay. See
// https://cloud.google.com/transcoder/docs/how-to/create-overlays#create-static-overlay
// for more information.
func createJobWithStaticOverlay(w io.Writer, projectID string, location string, inputURI string, overlayImageURI string, outputURI string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// inputURI := "gs://my-bucket/my-video-file"
	// overlayImageURI := "gs://my-bucket/my-overlay-image-file" - Must be a JPEG
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
									Codec:        "h264",
									BitrateBps:   550000,
									FrameRate:    60,
									HeightPixels: 360,
									WidthPixels:  640,
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
					Overlays: []*transcoderpb.Overlay{
						&transcoderpb.Overlay{
							Image: &transcoderpb.Overlay_Image{
								Uri: overlayImageURI,
								Resolution: &transcoderpb.Overlay_NormalizedCoordinate{
									X: 1,
									Y: 0.5,
								},
								Alpha: 1,
							},
							Animations: []*transcoderpb.Overlay_Animation{
								&transcoderpb.Overlay_Animation{
									AnimationType: &transcoderpb.Overlay_Animation_AnimationStatic{
										AnimationStatic: &transcoderpb.Overlay_AnimationStatic{
											Xy: &transcoderpb.Overlay_NormalizedCoordinate{
												X: 0,
												Y: 0,
											},
											StartTimeOffset: &duration.Duration{
												Seconds: 0,
											},
										},
									},
								},

								&transcoderpb.Overlay_Animation{
									AnimationType: &transcoderpb.Overlay_Animation_AnimationEnd{
										AnimationEnd: &transcoderpb.Overlay_AnimationEnd{
											StartTimeOffset: &duration.Duration{
												Seconds: 10,
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
		return fmt.Errorf("createJobWithStaticOverlay: %v", err)
	}

	fmt.Fprintf(w, "Job: %v", response.GetName())
	return nil
}

// [END transcoder_create_job_with_static_overlay]
