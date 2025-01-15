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

// [START speech_transcribe_model_selection_gcs]

import (
	"context"
	"fmt"
	"io"
	"strings"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
)

// transcribe_model_selection_gcs Transcribes the given audio file asynchronously with
// the selected model.
func transcribe_model_selection_gcs(w io.Writer) error {
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	audio := &speechpb.RecognitionAudio{
		AudioSource: &speechpb.RecognitionAudio_Uri{Uri: "gs://cloud-samples-tests/speech/Google_Gnome.wav"},
	}

	// The speech recognition model to use
	// See, https://cloud.google.com/speech-to-text/docs/speech-to-text-requests#select-model
	recognitionConfig := &speechpb.RecognitionConfig{
		Encoding:        speechpb.RecognitionConfig_LINEAR16,
		SampleRateHertz: 16000,
		LanguageCode:    "en-US",
		Model:           "video",
	}

	longRunningRecognizeRequest := &speechpb.LongRunningRecognizeRequest{
		Config: recognitionConfig,
		Audio:  audio,
	}

	operation, err := client.LongRunningRecognize(ctx, longRunningRecognizeRequest)
	if err != nil {
		return fmt.Errorf("error running recognize %w", err)
	}

	response, err := operation.Wait(ctx)
	if err != nil {
		return err
	}
	for i, result := range response.Results {
		alternative := result.Alternatives[0]
		fmt.Fprintf(w, "%s\n", strings.Repeat("-", 20))
		fmt.Fprintf(w, "First alternative of result %d", i)
		fmt.Fprintf(w, "Transcript: %s", alternative.Transcript)
	}
	return nil
}

// [END speech_transcribe_model_selection_gcs]
