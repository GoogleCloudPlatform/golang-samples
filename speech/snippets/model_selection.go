// Copyright 2019 Google LLC
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

// Package snippets contains speech examples.
package snippets

import (
	"fmt"
	"io"
	"os"
	"strings"

	// [START speech_transcribe_model_selection_imports]
	"context"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	// [END speech_transcribe_model_selection_imports]
)

// [START speech_transcribe_model_selection]

func modelSelection(w io.Writer) error {
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	data, err := os.ReadFile("../testdata/Google_Gnome.wav")
	if err != nil {
		return fmt.Errorf("ReadFile: %w", err)
	}

	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
			Model:           "video",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	}

	resp, err := client.Recognize(ctx, req)
	if err != nil {
		return fmt.Errorf("recognize: %w", err)
	}

	for i, result := range resp.Results {
		fmt.Fprintf(w, "%s\n", strings.Repeat("-", 20))
		fmt.Fprintf(w, "Result %d\n", i+1)
		for j, alternative := range result.Alternatives {
			fmt.Fprintf(w, "Alternative %d: %s\n", j+1, alternative.Transcript)
		}
	}
	return nil
}

// [END speech_transcribe_model_selection]
