// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command captionasync sends audio data to the Google Speech API
// and prints its transcript.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

const usage = `Usage: captionasync <audiofile>

Audio file must be a 16-bit signed little-endian encoded
with a sample rate of 16000.

The path to the audio file may be a GCS URI (gs://...).
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	var sendFunc func(*speech.Client, string) (*speechpb.LongRunningRecognizeResponse, error)

	path := os.Args[1]
	if strings.Contains(path, "://") {
		sendFunc = sendGCS
	} else {
		sendFunc = send
	}

	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := sendFunc(client, os.Args[1])
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

func send(client *speech.Client, filename string) (*speechpb.LongRunningRecognizeResponse, error) {
	ctx := context.Background()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Send the contents of the audio file with the encoding and
	// and sample rate information to be transcripted.
	req := &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	}

	op, err := client.LongRunningRecognize(ctx, req)
	if err != nil {
		return nil, err
	}
	return op.Wait(ctx)
}

func sendGCS(client *speech.Client, gcsURI string) (*speechpb.LongRunningRecognizeResponse, error) {
	ctx := context.Background()

	// Send the contents of the audio file with the encoding and
	// and sample rate information to be transcripted.
	req := &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: gcsURI},
		},
	}

	op, err := client.LongRunningRecognize(ctx, req)
	if err != nil {
		return nil, err
	}
	return op.Wait(ctx)
}
