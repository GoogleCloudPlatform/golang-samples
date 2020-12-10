// Copyright 2020 Google LLC
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

// [START transcoder_create_job_template]
import (
	"context"
	"fmt"
	"io"

	transcoder "cloud.google.com/go/video/transcoder/apiv1beta1"
	transcoderpb "google.golang.org/genproto/googleapis/cloud/video/transcoder/v1beta1"
)

// createJobTemplate creates a template for a job. See
// https://cloud.google.com/transcoder/docs/how-to/job-templates#create_job_templates
// for more information.
func createJobTemplate(w io.Writer, projectID string, location string, templateID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// templateID := "my-job-template"
	ctx := context.Background()
	client, err := transcoder.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	req := &transcoderpb.CreateJobTemplateRequest{
		Parent:        fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		JobTemplateId: templateID,
		JobTemplate: &transcoderpb.JobTemplate{
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
						Key: "video_stream1",
						ElementaryStream: &transcoderpb.ElementaryStream_VideoStream{
							VideoStream: &transcoderpb.VideoStream{
								Codec:        "h264",
								BitrateBps:   2500000,
								FrameRate:    60,
								HeightPixels: 720,
								WidthPixels:  1280,
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
					&transcoderpb.MuxStream{
						Key:               "hd",
						Container:         "mp4",
						ElementaryStreams: []string{"video_stream1", "audio_stream0"},
					},
				},
			},
		},
	}

	response, err := client.CreateJobTemplate(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateJobTemplate: %v", err)
	}

	fmt.Fprintf(w, "Job template: %v", response.GetName())
	return nil
}

// [END transcoder_create_job_template]
