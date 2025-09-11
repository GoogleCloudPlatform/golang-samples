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

// [START speech_transcribe_diarization_gcs_beta]

import (
	"context"
	"fmt"
	"io"
	"strings"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
)

// transcribe_diarization_gcs_beta Transcribes a remote audio file using speaker diarization.
func transcribe_diarization_gcs_beta(w io.Writer) error {
	// Google Cloud Storage URI pointing to the audio content.

	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	diarizationConfig := &speechpb.SpeakerDiarizationConfig{
		EnableSpeakerDiarization: true,
		MinSpeakerCount:          2,
		MaxSpeakerCount:          2,
	}

	recognitionConfig := &speechpb.RecognitionConfig{
		Encoding:          speechpb.RecognitionConfig_LINEAR16,
		SampleRateHertz:   8000,
		LanguageCode:      "en-US",
		DiarizationConfig: diarizationConfig,
	}

	// Set the remote path for the audio file
	audio := &speechpb.RecognitionAudio{
		AudioSource: &speechpb.RecognitionAudio_Uri{Uri: "gs://cloud-samples-tests/speech/commercial_mono.wav"},
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

	// Speaker Tags are only included in the last result object, which has only one
	// alternative.
	alternative := response.Results[len(response.Results)-1].Alternatives[0]

	wordInfo := alternative.GetWords()[0]
	currentSpeakerTag := wordInfo.GetSpeakerTag()

	var speakerWords strings.Builder

	speakerWords.WriteString(fmt.Sprintf("Speaker %d: %s", wordInfo.GetSpeakerTag(), wordInfo.GetWord()))

	// For each word, get all the words associated with one speaker, once the speaker changes,
	// add a new line with the new speaker and their spoken words.
	for i := 1; i < len(alternative.Words); i++ {
		wordInfo := alternative.Words[i]
		if currentSpeakerTag == wordInfo.GetSpeakerTag() {
			speakerWords.WriteString(" ")
			speakerWords.WriteString(wordInfo.GetWord())
		} else {
			speakerWords.WriteString(fmt.Sprintf("\nSpeaker %d: %s",
				wordInfo.GetSpeakerTag(), wordInfo.GetWord()))
			currentSpeakerTag = wordInfo.GetSpeakerTag()
		}
	}
	fmt.Fprint(w, speakerWords.String())
	return nil
}

// [END speech_transcribe_diarization_gcs_beta]
