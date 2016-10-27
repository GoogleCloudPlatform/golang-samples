// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command captionasync sends audio data to the Google Speech API
// and pulls the operation status and the transcript.
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gogo/protobuf/proto"

	"golang.org/x/net/context"

	speech "cloud.google.com/go/speech/apiv1beta1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1beta1"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
)

const usage = `Usage: captionasync <audiofile>

Audio file is required to be 16-bit signed little-endian encoded
with a sample rate of 16000.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	opName, err := send(client, os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	resp, err := wait(client, opName)
	if err != nil {
		log.Fatal(err)
	}

	// Print the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}
}

func send(client *speech.Client, filename string) (string, error) {
	ctx := context.Background()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// Send the contents of the audio file with the encoding and
	// and sample rate information to be transcripted.
	req := &speechpb.AsyncRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:   speechpb.RecognitionConfig_LINEAR16,
			SampleRate: 16000,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	}

	op, err := client.AsyncRecognize(ctx, req)
	if err != nil {
		return "", err
	}
	return op.Name(), nil
}

func wait(client *speech.Client, opName string) (*speechpb.AsyncRecognizeResponse, error) {
	ctx := context.Background()

	opClient := longrunningpb.NewOperationsClient(client.Connection())
	var op *longrunningpb.Operation
	var err error
	for {
		op, err = opClient.GetOperation(ctx, &longrunningpb.GetOperationRequest{
			Name: opName,
		})
		if err != nil {
			return nil, err
		}
		if op.Done {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	switch {
	case op.GetError() != nil:
		return nil, fmt.Errorf("recieved error in response: %v", op.GetError())
	case op.GetResponse() != nil:
		var resp speechpb.AsyncRecognizeResponse
		if err := proto.Unmarshal(op.GetResponse().Value, &resp); err != nil {
			return nil, err
		}
		return &resp, nil
	}

	// should never happen.
	return nil, errors.New("no response")
}
