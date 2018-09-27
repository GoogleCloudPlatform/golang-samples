// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package snippets contains speech examples.
package snippets

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	// [START imports]
	"context"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	// [END imports]
)

// [START speech_transcribe_auto_punctuation]

func autoPunctuation(w io.Writer, path string) error {
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}

	// path = "../testdata/commercial_mono.wav"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("ReadFile: %v", err)
	}

	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 8000,
			LanguageCode:    "en-US",
			// Enable automatic punctuation.
			EnableAutomaticPunctuation: true,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	if err != nil {
		return fmt.Errorf("Recognize: %v", err)
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

// [END speech_transcribe_auto_punctuation]
