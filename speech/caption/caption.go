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

	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	speech "google.golang.org/genproto/googleapis/cloud/speech/v1beta1"
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

	ctx := context.Background()
	conn, err := transport.DialGRPC(ctx,
		option.WithEndpoint("speech.googleapis.com:443"),
		option.WithScopes("https://www.googleapis.com/auth/cloud-platform"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := speech.NewSpeechClient(conn)
	// TODO(jbd): switch to the bidirectional streaming api
	// and send data in small chunks.
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Send the contents of the audio file with the encoding and
	// and sample rate information to be transcripted.
	rresp, err := recognize(ctx, c, &data)
	if err != nil {
		log.Fatal(err)
	}

	// Print the results.
	for _, result := range rresp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}
}

func recognize(ctx context.Context, client speech.SpeechClient, data *[]byte) (*speech.SyncRecognizeResponse, error) {
	return client.SyncRecognize(ctx, &speech.SyncRecognizeRequest{
		Config: &speech.RecognitionConfig{
			Encoding:   speech.RecognitionConfig_LINEAR16,
			SampleRate: 16000,
		},
		Audio: &speech.RecognitionAudio{
			AudioSource: &speech.RecognitionAudio_Content{Content: *data},
		},
	})
}
