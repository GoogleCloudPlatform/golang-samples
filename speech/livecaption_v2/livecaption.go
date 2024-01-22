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

// Command livecaption pipes the stdin audio data to
// Google Speech API and outputs the transcript.
//
// As an example, gst-launch can be used to capture the mic input:
//
//	$ gst-launch-1.0 -v pulsesrc ! audioconvert ! audioresample ! audio/x-raw,channels=1,rate=16000 ! filesink location=/dev/stdout | livecaption <project_id>

package main

// [START speech_transcribe_streaming_mic]
import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	speech "cloud.google.com/go/speech/apiv2"
	"cloud.google.com/go/speech/apiv2/speechpb"
)

var projectID string

const location = "global"

func main() {
	ctx := context.Background()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <Project_id>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "<projectID> must be a project_id to a valid gcp projectID with speech api enabled.\n")

	}
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("Please pass the project_id as a command line argument. Should be a valid project_id with stt api enabled.")
	}
	projectID = flag.Arg(0)

	if projectID == "" {
		log.Fatalf("Project is is required parameter: %s", projectID)
	}

	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	stream, err := client.StreamingRecognize(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := stream.Send(&speechpb.StreamingRecognizeRequest{
		Recognizer: fmt.Sprintf("projects/%s/locations/%s/recognizers/_", projectID, location),
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					// In case of specific file encoding , so specify the decoding config.
					//DecodingConfig: &speechpb.RecognitionConfig_AutoDecodingConfig{},
					DecodingConfig: &speechpb.RecognitionConfig_ExplicitDecodingConfig{
						ExplicitDecodingConfig: &speechpb.ExplicitDecodingConfig{
							Encoding:          speechpb.ExplicitDecodingConfig_LINEAR16,
							SampleRateHertz:   16000,
							AudioChannelCount: 1,
						},
					},
					Model:         "long",
					LanguageCodes: []string{"en-US"},
					Features: &speechpb.RecognitionFeatures{
						MaxAlternatives: 2,
					},
				},
				StreamingFeatures: &speechpb.StreamingRecognitionFeatures{InterimResults: true},
			},
		},
	}); err != nil {
		log.Fatal(err)
	}

	go func() {
		// Pipe stdin to the API.
		buf := make([]byte, 1024)

		for {

			n, err := os.Stdin.Read(buf)

			if n > 0 {
				if err := stream.Send(&speechpb.StreamingRecognizeRequest{
					Recognizer: fmt.Sprintf("projects/%s/locations/%s/recognizers/_", projectID, location),
					StreamingRequest: &speechpb.StreamingRecognizeRequest_Audio{
						Audio: buf[:n],
					},
				}); err != nil {
					log.Printf("Could not send audio: %v", err)
				}
			}
			if err == io.EOF {
				// Nothing else to pipe, close the stream.
				if err := stream.CloseSend(); err != nil {
					log.Fatalf("Could not close stream: %v", err)
				}
				return
			}
			if err != nil {
				log.Printf("Could not read from stdin: %v", err)
				continue
			}
		}
	}()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Printf("EOF break")
			break
		}
		if err != nil {
			log.Fatalf("Could not recognize: %v", err)
		} else {
			// It seems like the new response api does not have a field called Error
			for _, result := range resp.Results {
				//fmt.Printf("Result: %+v\n", result)
				if len(result.Alternatives) > 0 {
					if result.IsFinal == true {
						log.Println("result", result.Alternatives[0].Transcript, result.IsFinal)
					}

				}
			}
		}

	}
}

// [END speech_transcribe_streaming_mic]
