# Copyright 2024 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


package main

import (
	"context"
	"log"
	"os"
	"regexp"
	"testing"

	"cloud.google.com/go/speech/apiv2”
	"google.golang.org/api/option"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v2”
)

var resourcesDir = "resources"

func TestTranscribeStreamingV2(t *testing.T) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		t.Fatalf("GOOGLE_CLOUD_PROJECT not set")
	}

	ctx := context.Background()

	client, err := speech.NewClient(ctx, option.WithCredentialsFile("path_to_your_credentials.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	filePath := filepath.Join(resourcesDir, "audio.wav")

	audioData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read audio file: %v", err)
	}

	streamingReq := &speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					Encoding:        speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz: 16000,
					LanguageCode:    "en-US",
					Model:           "long-audio",
				},
			},
		},
	}

	stream, err := client.StreamingRecognize(ctx)
	if err != nil {
		t.Fatalf("Failed to create streaming request: %v", err)
	}
	defer stream.CloseSend()

	if err := stream.Send(streamingReq); err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Send audio data
	for i := 0; i < len(audioData); i += 32000 {
		end := i + 32000
		if end > len(audioData) {
			end = len(audioData)
		}
		req := &speechpb.StreamingRecognizeRequest{
			StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
				AudioContent: audioData[i:end],
			},
		}
		if err := stream.Send(req); err != nil {
			t.Fatalf("Failed to send audio content: %v", err)
		}
	}

	resp, err := stream.Recv()
	if err != nil {
		t.Fatalf("Failed to receive response: %v", err)
	}

	var transcript string
	for _, result := range resp.GetResults() {
		for _, alt := range result.GetAlternatives() {
			transcript += alt.GetTranscript()
		}
	}

	match, err := regexp.MatchString(`how old is the Brooklyn Bridge`, transcript)
	if err != nil {
		t.Fatalf("Failed to match transcript: %v", err)
	}

	if !match {
		t.Errorf("Expected transcript not found in response: %s", transcript)
	}
}

