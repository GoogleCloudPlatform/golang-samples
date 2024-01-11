// Copyright 2023 Google LLC
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

package snippets

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	speech "cloud.google.com/go/speech/apiv2"
	"cloud.google.com/go/speech/apiv2/speechpb"
)

func transcribeStreamingV2(w io.Writer, projectId, path string) error {
	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	recognizer := fmt.Sprintf("projects/%s/locations/global/recognizers/_", projectId)

	if err != nil {
		return fmt.Errorf("speech.NewClient: %w", err)
	}
	defer client.Close()

	// path = "../testdata/commercial_mono.wav"
	content, err := os.ReadFile(path)

	if err != nil {
		return fmt.Errorf("file: %s, ReadFile Error: %w", path, err)
	}

	//chunk the audio data to simulate streaming
	chunkLength := len(content) / 500
	audioRequests := make([]speechpb.StreamingRecognizeRequest_Audio, 0, len(content)/chunkLength)
	for i := 0; i < len(content); i += chunkLength {
		end := i + chunkLength
		if end > len(content) {
			end = len(content)
		}
		audioRequests = append(audioRequests, speechpb.StreamingRecognizeRequest_Audio{Audio: content[i:end]})
	}

	languageCodes := []string{"en-US"}
	recognitionConfig := &speechpb.RecognitionConfig{
		DecodingConfig: &speechpb.RecognitionConfig_AutoDecodingConfig{},
		Model:          "short",
		LanguageCodes:  languageCodes,
	}

	streamingConfig := &speechpb.StreamingRecognizeRequest_StreamingConfig{
		StreamingConfig: &speechpb.StreamingRecognitionConfig{Config: recognitionConfig},
	}
	configRequest := &speechpb.StreamingRecognizeRequest{
		Recognizer:       recognizer,
		StreamingRequest: streamingConfig,
	}

	stream, err := client.StreamingRecognize(ctx)

	go func() {
		stream.Send(configRequest) // First message sent should be the recognition config
		for _, request := range audioRequests {
			stream.Send(&speechpb.StreamingRecognizeRequest{
				Recognizer:       recognizer,
				StreamingRequest: &request,
			})
		}
		stream.CloseSend()
	}()

	for {
		response, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("file: %s, StreamingRecognize: %w", path, err)
		}

		for i, result := range response.Results {
			fmt.Fprintf(w, "%s\n", strings.Repeat("-", 20))
			fmt.Fprintf(w, "Result %d\n", i+1)
			for j, alternative := range result.Alternatives {
				fmt.Fprintf(w, "Alternative %d: %s\n", j+1, alternative.Transcript)
			}
		}
	}

	return nil
}
