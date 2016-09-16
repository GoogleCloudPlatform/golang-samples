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

	// [START imports]
	"golang.org/x/net/context"

	speech "cloud.google.com/go/speech/apiv1beta1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1beta1"
	// [END imports]
)

const usage = `Usage: caption <audiofile>

Audio file is required to be 16-bit signed little-endian encoded
with a sample rate of 16000.
`

func main() {
	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	// Perform the request.
	resp, err := recognize(os.Args[1])
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

func recognize(file string) (*speechpb.SyncRecognizeResponse, error) {
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
	resp, err := client.SyncRecognize(ctx, &speechpb.SyncRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:   speechpb.RecognitionConfig_LINEAR16,
			SampleRate: 16000,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	// [END request]
	return resp, err
}
