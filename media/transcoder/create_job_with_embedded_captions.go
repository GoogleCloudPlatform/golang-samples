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

// [START transcoder_create_job_with_embedded_captions]
import (
	"context"
	"fmt"
	"io"

	transcoder "cloud.google.com/go/video/transcoder/apiv1"
	"cloud.google.com/go/video/transcoder/apiv1/transcoderpb"
)

// createJobWithEmbeddedCaptions creates a job that embeds closed captions in the
// output video. See https://cloud.google.com/transcoder/docs/how-to/captions-and-subtitles
// for more information.
func createJobWithEmbeddedCaptions(w io.Writer, projectID string, location string, inputVideoURI string, inputCaptionsURI string, outputURI string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// inputVideoURI := "gs://my-bucket/my-video-file"
	// inputCaptionsURI := "gs://my-bucket/my-captions-file"
	// outputURI := "gs://my-bucket/my-output-folder/"

	ctx := context.Background()
	client, err := transcoder.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	// Set up elementary streams. The InputKey field refers to inputs in
	// the Inputs array defined the job config.
	elementaryStreams := []*transcoderpb.ElementaryStream{
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
		{
			Key: "cea_stream0",
			ElementaryStream: &transcoderpb.ElementaryStream_TextStream{
				TextStream: &transcoderpb.TextStream{
					Codec: "cea608",
					Mapping: []*transcoderpb.TextStream_TextMapping{
						{
							AtomKey:    "atom0",
							InputKey:   "caption_input0",
							InputTrack: 0,
						},
					},
					LanguageCode: "en-US",
					DisplayName:  "English",
				},
			},
		},
	}

	req := &transcoderpb.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Job: &transcoderpb.Job{
			OutputUri: outputURI,
			JobConfig: &transcoderpb.Job_Config{
				Config: &transcoderpb.JobConfig{
					Inputs: []*transcoderpb.Input{
						{
							Key: "input0",
							Uri: inputVideoURI,
						},
						{
							Key: "caption_input0",
							Uri: inputCaptionsURI,
						},
					},
					EditList: []*transcoderpb.EditAtom{
						{
							Key:    "atom0",
							Inputs: []string{"input0", "caption_input0"},
						},
					},
					ElementaryStreams: elementaryStreams,
					MuxStreams: []*transcoderpb.MuxStream{
						{
							Key:               "sd-hls",
							Container:         "ts",
							ElementaryStreams: []string{"video_stream0", "audio_stream0"},
						},
						{
							Key:               "sd-dash",
							Container:         "fmp4",
							ElementaryStreams: []string{"video_stream0"},
						},
						{
							Key:               "audio-dash",
							Container:         "fmp4",
							ElementaryStreams: []string{"audio_stream0"},
						},
					},
					Manifests: []*transcoderpb.Manifest{
						{
							FileName:   "manifest.m3u8",
							Type:       transcoderpb.Manifest_HLS,
							MuxStreams: []string{"sd-hls"},
						},
						{
							FileName:   "manifest.mpd",
							Type:       transcoderpb.Manifest_DASH,
							MuxStreams: []string{"sd-dash", "audio-dash"},
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

// [END transcoder_create_job_with_embedded_captions]
