// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command auto_punctuation sends audio data to the Google Speech API
// and prints its transcript with auto punctuation enabled.
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	// [START imports]
	"golang.org/x/net/context"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	// [END imports]
)

const usage = `Usage: auto_punctuation <audiofile>

Audio file must be a 16-bit signed little-endian encoded
with a sample rate of 8000.
`

func main() {
	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	err := autoPunctuation(os.Stdout, os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

func autoPunctuation(w io.Writer, path string) error {
	// [START speech_transcribe_auto_punctuation]
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// path := "../testdata/commercial_mono.wav"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 8000,
			LanguageCode:    "en-US",
			// Enable automatic punctuation
			EnableAutomaticPunctuation: true,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})

	for i, result := range resp.Results {
		fmt.Fprintf(w, "%s\n", strings.Repeat("-", 20))
		fmt.Fprintf(w, "Result %d\n", i + 1)
		for j, alternative := range result.Alternatives {
			fmt.Fprintf(w, "Alternative %d: %s\n", j + 1, alternative.Transcript)
		}
	}
	// [END speech_transcribe_auto_punctuation]
	return err
}
