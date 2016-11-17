// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command livecaption pipes the stdin audio data to
// Google Speech API and outputs the transcript.
//
// As an example, gst-launch can be used to capture the mic input:
//
//    $ gst-launch-1.0 -v pulsesrc ! audioconvert ! audioresample ! audio/x-raw,channels=1,rate=16000 ! filesink location=/dev/stdout | livecaption
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	speech "google.golang.org/genproto/googleapis/cloud/speech/v1beta1"
)

func main() {
	ctx := context.Background()
	conn, err := transport.DialGRPC(ctx,
		option.WithEndpoint("speech.googleapis.com:443"),
		option.WithScopes("https://www.googleapis.com/auth/cloud-platform"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// [START speech_streaming_mic_recognize]
	stream, err := speech.NewSpeechClient(conn).StreamingRecognize(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// send the initial configuration message.
	if err := stream.Send(&speech.StreamingRecognizeRequest{
		StreamingRequest: &speech.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speech.StreamingRecognitionConfig{
				Config: &speech.RecognitionConfig{
					Encoding:   speech.RecognitionConfig_LINEAR16,
					SampleRate: 16000,
				},
			},
		},
	}); err != nil {
		log.Fatal(err)
	}

	go func() {
		// pipe stdin to the API
		buf := make([]byte, 1024)
		for {
			n, err := os.Stdin.Read(buf)
			if err == io.EOF {
				return // nothing else to pipe, kill this goroutine
			}
			if err != nil {
				log.Printf("Could not read from stdin: %v", err)
				continue
			}
			if err = stream.Send(&speech.StreamingRecognizeRequest{
				StreamingRequest: &speech.StreamingRecognizeRequest_AudioContent{
					AudioContent: buf[:n],
				},
			}); err != nil {
				log.Printf("Could not send audio: %v", err)
			}
		}
	}()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot stream results: %v", err)
		}
		if err := resp.Error; err != nil {
			log.Fatalf("Could not recognize: %v", err)
		}
		for _, result := range resp.Results {
			fmt.Printf("Result: %+v\n", result)
		}
	}
	// [END speech_streaming_mic_recognize]
}
