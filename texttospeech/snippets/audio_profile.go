// Copyright 2019 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
