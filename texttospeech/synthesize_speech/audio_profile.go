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

// [START tts_synthesize_text_audio_profile]
// [START texttospeech_audio_profile]

import (
	"fmt"
	"io"
	"os"

	"context"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

// audioProfile generates audio from text using a custom synthesizer like a telephone call.
func audioProfile(w io.Writer, text string, outputFile string) error {
	// text := "hello"
	// outputFile := "out.mp3"

	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{LanguageCode: "en-US"},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding:    texttospeechpb.AudioEncoding_MP3,
			EffectsProfileId: []string{"telephony-class-application"},
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, req)
	if err != nil {
		return fmt.Errorf("SynthesizeSpeech: %w", err)
	}

	if err = os.WriteFile(outputFile, resp.AudioContent, 0644); err != nil {
		return err
	}

	fmt.Fprintf(w, "Audio content written to file: %v\n", outputFile)

	return nil
}

// [END texttospeech_audio_profile]
// [END tts_synthesize_text_audio_profile]
