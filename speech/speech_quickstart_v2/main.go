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
	"cloud.google.com/go/speech/apiv2/speechpb"
)

func main() {

	//Creates a context
	ctx := context.Background()

	//Create a new client
	client, err := speech.NewClient(ctx)

	//Check for error if any
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// The path to the remote audio file to transcribe.
	fileURI := "gs://cloud-samples-data/speech/brooklyn_bridge.raw"

	// Detects speech in the audio file.
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			DecodingConfig: &speechpb.RecognitionConfig_AutoDecodingConfig{},
			LanguageCodes:  []string{"en-US"},
		},
		AudioSource: &speechpb.RecognizeRequest_Uri{Uri: fileURI},
	})
	if err != nil {
		log.Fatalf("failed to recognize: %v", err)
	}

	// Prints the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}

}

// [END speech_quickstart_v2]
