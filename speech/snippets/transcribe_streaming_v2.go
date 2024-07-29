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

// Command livecaption_from_file streams a local audio file to
// Google Speech API and outputs the transcript.

package snippets

// [START speech_transcribe_streaming_v2]
import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	speech "cloud.google.com/go/speech/apiv2"
	"cloud.google.com/go/speech/apiv2/speechpb"
)

func transcribeStreamingV2(w io.Writer) error {
	path := "../resources/commercial_mono.wav"
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	const location = "global"

	audioFile, err := filepath.Abs(path)
	if err != nil {
		log.Println("Failed to load file: ", path)
		return err
	}
	f, err := os.Open(audioFile)
	if err != nil {
		return err
	}
	defer f.Close()

	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	stream, err := client.StreamingRecognize(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	// Send the initial configuration message.
	err = stream.Send(&speechpb.StreamingRecognizeRequest{
		Recognizer: fmt.Sprintf("projects/%s/locations/%s/recognizers/_", projectID, location),
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					// In case of specific file encoding , so specify the decoding config.
					DecodingConfig: &speechpb.RecognitionConfig_AutoDecodingConfig{},
					Model:          "long",
					LanguageCodes:  []string{"en-US"},
					Features: &speechpb.RecognitionFeatures{
						MaxAlternatives: 2,
					},
				},
				StreamingFeatures: &speechpb.StreamingRecognitionFeatures{InterimResults: true},
			},
		},
	})
	if err != nil {
		return err
	}

	go func() error {
		buf := make([]byte, 1024)
		for {
			n, err := f.Read(buf)
			if n > 0 {
				if err := stream.Send(&speechpb.StreamingRecognizeRequest{
					Recognizer: fmt.Sprintf("projects/%s/locations/%s/recognizers/_", projectID, location),
					StreamingRequest: &speechpb.StreamingRecognizeRequest_Audio{
						Audio: buf[:n],
					},
				}); err != nil {
					return fmt.Errorf("could not send audio: %w", err)
				}
			}
			if err == io.EOF {
				// Nothing else to pipe, close the stream.
				if err := stream.CloseSend(); err != nil {
					return fmt.Errorf("could not close stream: %w", err)
				}
				return nil
			}
			if err != nil {
				log.Printf("Could not read from %s: %v", audioFile, err)
				continue
			}
		}
	}()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("cannot stream results: %w", err)
		}
		for i, result := range resp.Results {
			fmt.Fprintf(w, "%s\n", strings.Repeat("-", 20))
			fmt.Fprintf(w, "Result %d\n", i+1)
			for j, alternative := range result.Alternatives {
				fmt.Fprintf(w, "Alternative %d is_final: %t : %s\n", j+1, result.IsFinal, alternative.Transcript)
			}
		}
	}
	return nil
}

// [END speech_transcribe_streaming_v2]
