# Copyright 2024 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/speech/apiv2”
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v2”
)

func main() {
	// Parse command-line arguments
	flag.Parse()
	args := flag.Args()

	// Get the project ID and audio file path from command-line arguments
	projectID := args[0]
	audioFile := args[1]

	// Create a context
	ctx := context.Background()

	// Create a client
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Open the audio file
	audioData, err := os.ReadFile(audioFile)
	if err != nil {
		log.Fatalf("Failed to read audio file: %v", err)
	}

	// Split audio data into chunks
	var audioChunks [][]byte
	scanner := bufio.NewScanner(bytes.NewReader(audioData))
	for scanner.Scan() {
		audioChunks = append(audioChunks, scanner.Bytes())
	}

	// Create the streaming request
	streamingReqs := make([]*speechpb.StreamingRecognizeRequest, len(audioChunks))
	recognitionConfig := &speechpb.RecognitionConfig{
		Encoding:        speechpb.RecognitionConfig_LINEAR16,
		SampleRateHertz: 16000,
		LanguageCode:    "en-US",
		Model:           "long",
	}
	streamingConfig := &speechpb.StreamingRecognitionConfig{
		Config: recognitionConfig,
	}
	for i, chunk := range audioChunks {
		streamingReqs[i] = &speechpb.StreamingRecognizeRequest{
			StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
				AudioContent: chunk,
			},
		}
	}

	// Create the streaming request
	stream, err := client.StreamingRecognize(ctx)
	if err != nil {
		log.Fatalf("Failed to create streaming request: %v", err)
	}

	// Send the streaming request
	for _, req := range streamingReqs {
		if err := stream.Send(req); err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}
	}

	// Receive and print responses
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive response: %v", err)
		}
		for _, result := range resp.GetResults() {
			for _, alt := range result.GetAlternatives() {
				fmt.Printf("Transcript: %s\n", alt.GetTranscript())
			}
		}
	}
}

