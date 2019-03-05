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
	"io/ioutil"

	// [START imports]
	"context"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	// [END imports]
)

// [START speech_transcribe_auto_punctuation]

func audioProfile(w io.Writer, text string, outputFile string) error {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}

	// text = "hello"
	resp, err := client.SynthesizeSpeech(ctx, &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{LanguageCode: "en-US"},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding:    texttospeechpb.AudioEncoding_MP3,
			EffectsProfileId: []string{"telephony-class-application"},
		},
	})
	if err != nil {
		return fmt.Errorf("SynthesizeSpeech: %v", err)
	}

	// outputFile = "output.mp3"
	err = ioutil.WriteFile(outputFile, resp.AudioContent, 0644)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Audio content written to file: %v\n", outputFile)
	return nil
}

// [END speech_transcribe_auto_punctuation]
