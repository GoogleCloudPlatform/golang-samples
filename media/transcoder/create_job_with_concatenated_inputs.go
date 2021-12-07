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

// [START transcoder_create_job_with_concatenated_inputs]
import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/golang/protobuf/ptypes/duration"

	transcoder "cloud.google.com/go/video/transcoder/apiv1"
	transcoderpb "google.golang.org/genproto/googleapis/cloud/video/transcoder/v1"
)

// createJobWithConcatenatedInputs creates a job that concatenates two input
// videos. See https://cloud.google.com/transcoder/docs/how-to/concatenate-videos
// for more information.
func createJobWithConcatenatedInputs(w io.Writer, projectID string, location string, input1URI string, startTimeOffset1 float64, endTimeOffset1 float64, input2URI string, startTimeOffset2 float64, endTimeOffset2 float64, outputURI string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// input1URI := "gs://my-bucket/my-video-file1"
	// startTimeOffset1 := 0
	// endTimeOffset1 := 8.1
	// input2URI := "gs://my-bucket/my-video-file2"
	// startTimeOffset2 := 3.5
	// endTimeOffset2 := 15
	// outputURI := "gs://my-bucket/my-output-folder/"

	whole, frac := math.Modf(startTimeOffset1)
	frac *= 1000000000
	var startTimeOffset1NanoSec int32 = int32(frac)
	var startTimeOffset1Sec int64 = int64(whole)

	whole, frac = math.Modf(endTimeOffset1)
	frac *= 1000000000
	var endTimeOffset1NanoSec int32 = int32(frac)
	var endTimeOffset1Sec int64 = int64(whole)

	whole, frac = math.Modf(startTimeOffset2)
	frac *= 1000000000
	var startTimeOffset2NanoSec int32 = int32(frac)
	var startTimeOffset2Sec int64 = int64(whole)

	whole, frac = math.Modf(endTimeOffset2)
	frac *= 1000000000
	var endTimeOffset2NanoSec int32 = int32(frac)
	var endTimeOffset2Sec int64 = int64(whole)

	ctx := context.Background()
	client, err := transcoder.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	req := &transcoderpb.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Job: &transcoderpb.Job{
			OutputUri: outputURI,
			JobConfig: &transcoderpb.Job_Config{
				Config: &transcoderpb.JobConfig{
					Inputs: []*transcoderpb.Input{
						&transcoderpb.Input{
							Key: "input1",
							Uri: input1URI,
						},
						&transcoderpb.Input{
							Key: "input2",
							Uri: input2URI,
						},
					},
					EditList: []*transcoderpb.EditAtom{
						&transcoderpb.EditAtom{
							Key:    "atom1",
							Inputs: []string{"input1"},
							StartTimeOffset: &duration.Duration{
								Seconds: startTimeOffset1Sec,
								Nanos:   startTimeOffset1NanoSec,
							},
							EndTimeOffset: &duration.Duration{
								Seconds: endTimeOffset1Sec,
								Nanos:   endTimeOffset1NanoSec,
							},
						},
						&transcoderpb.EditAtom{
							Key:    "atom2",
							Inputs: []string{"input2"},
							StartTimeOffset: &duration.Duration{
								Seconds: startTimeOffset2Sec,
								Nanos:   startTimeOffset2NanoSec,
							},
							EndTimeOffset: &duration.Duration{
								Seconds: endTimeOffset2Sec,
								Nanos:   endTimeOffset2NanoSec,
							},
						},
					},
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
				},
			},
		},
	}
	// Creates the job. Jobs take a variable amount of time to run.
	// You can query for the job state; see getJob() in get_job.go.
	response, err := client.CreateJob(ctx, req)
	if err != nil {
		return fmt.Errorf("createJobWithConcatenatedInputs: %v", err)
	}

	fmt.Fprintf(w, "Job: %v", response.GetName())
	return nil
}

// [END transcoder_create_job_with_concatenated_inputs]
