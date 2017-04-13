// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command caption reads an audio file and outputs the transcript for it.
package main

import (
	"flag"
	"fmt"
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

const usage = `Usage: caption <audiofile>

Audio file must be a 16-bit signed little-endian encoded
with a sample rate of 16000.

The path to the audio file may be a GCS URI (gs://...).
`

func main() {
	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	var runFunc func(string) (*speechpb.RecognizeResponse, error)

	path := os.Args[1]
	if strings.Contains(path, "://") {
		runFunc = recognizeGCS
	} else {
		runFunc = recognize
	}

	// Perform the request.
	resp, err := runFunc(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// [START print]
	// Print the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}
	// [END print]
}

func recognizeGCS(gcsURI string) (*speechpb.RecognizeResponse, error) {
	ctx := context.Background()

	// [START init_gcs]
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// [END init_gcs]

	// [START request_gcs]
	// Send the request with the URI (gs://...)
	// and sample rate information to be transcripted.
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: gcsURI},
		},
	})
	// [END request_gcs]
	return resp, err
}

func recognize(file string) (*speechpb.RecognizeResponse, error) {
	ctx := context.Background()

	// [START init]
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// [END init]

	// [START request]
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// Send the contents of the audio file with the encoding and
	// and sample rate information to be transcripted.
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	// [END request]
	return resp, err
}
