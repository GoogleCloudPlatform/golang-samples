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


// [START speech_quickstart_v2]

// Sample speech-quickstart_v2 uses the Google Cloud Speech API to transcribe audio.

package main
import (
	"context"
	"fmt"
	"log"
	
        speech "cloud.google.com/go/speech/apiv2"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v2"
)

func main() {
	
        // Parse command-line arguments
	projectID := flag.String("project_id", "", "GCP Project ID")
	flag.Parse()
        
        // The path to the remote audio file to transcribe.
        audioFile := "gs://cloud-samples-data/speech/brooklyn_bridge.raw"
	
        // Check if required arguments are provided
        if *projectID == "" || *audioFile == "" {
		flag.Usage()
		return
	}
	
        // Create a new context
	ctx := context.Background()
	
        // Create a client
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
        defer client.Close()
	
        // Detects speech in the audio file provided via command-line argument
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: *audioFile},
		},
	})
	if err != nil {
		log.Fatalf("failed to recognize: %v", err)
	}
	
        // Prints the results
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%f)\n", alt.Transcript, alt.Confidence)
		}
	}
}

// [END speech_quickstart_v2]
