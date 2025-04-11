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

package transcoder

// [START transcoder_create_job_with_concatenated_inputs]
import (
	"context"
	"fmt"
	"io"
	"time"

	transcoder "cloud.google.com/go/video/transcoder/apiv1"
	"cloud.google.com/go/video/transcoder/apiv1/transcoderpb"
	"google.golang.org/protobuf/types/known/durationpb"
)

// createJobWithConcatenatedInputs creates a job that concatenates two input
// videos. See https://cloud.google.com/transcoder/docs/how-to/concatenate-videos
// for more information.
func createJobWithConcatenatedInputs(w io.Writer, projectID string, location string, input1URI string, startTimeInput1 time.Duration, endTimeInput1 time.Duration, input2URI string, startTimeInput2 time.Duration, endTimeInput2 time.Duration, outputURI string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	//
	// input1URI := "gs://my-bucket/my-video-file1"
	// startTimeInput1 := 0*time.Second
	// endTimeInput1 := 8*time.Second + 100*time.Millisecond
	//
	// input2URI := "gs://my-bucket/my-video-file2"
	// startTimeInput2 := 3*time.Second + 500*time.Millisecond
	// endTimeInput2 := 15*time.Second
	//
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
			OutputUri: outputURI,
			JobConfig: &transcoderpb.Job_Config{
				Config: &transcoderpb.JobConfig{
					Inputs: []*transcoderpb.Input{
						{
							Key: "input1",
							Uri: input1URI,
						},
						{
							Key: "input2",
							Uri: input2URI,
						},
					},
					EditList: []*transcoderpb.EditAtom{
						{
							Key:             "atom1",
							Inputs:          []string{"input1"},
							StartTimeOffset: durationpb.New(startTimeInput1),
							EndTimeOffset:   durationpb.New(endTimeInput1),
						},
						{
							Key:             "atom2",
							Inputs:          []string{"input2"},
							StartTimeOffset: durationpb.New(startTimeInput2),
							EndTimeOffset:   durationpb.New(endTimeInput2),
						},
					},
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
				},
			},
		},
	}
	// Creates the job. Jobs take a variable amount of time to run.
	// You can query for the job state; see getJob() in get_job.go.
	response, err := client.CreateJob(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateJob: %w", err)
	}

	fmt.Fprintf(w, "Job: %v", response.GetName())
	return nil
}

// [END transcoder_create_job_with_concatenated_inputs]
