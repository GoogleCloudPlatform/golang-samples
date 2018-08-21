// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command wordoffset sends audio data to the Google Speech API
// and prints word offset information.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

const usage = `Usage: wordoffset <audiofile>

Audio file must be a 16-bit signed little-endian encoded
with a sample rate of 16000.

The path to the audio file may be a GCS URI (gs://...).
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	var sendFunc func(*speech.Client, io.Writer, string) error

	path := os.Args[1]
	if strings.Contains(path, "://") {
		sendFunc = asyncWords
	} else {
		sendFunc = syncWords
	}

	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := sendFunc(client, os.Stdout, os.Args[1]); err != nil {
		log.Fatal(err)
	}
}

// [START speech_transcribe_async_time_offsets_gcs]

func asyncWords(client *speech.Client, out io.Writer, gcsURI string) error {
	ctx := context.Background()

	// Send the contents of the audio file with the encoding and
	// and sample rate information to be transcripted.
	req := &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:              speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz:       16000,
			LanguageCode:          "en-US",
			EnableWordTimeOffsets: true,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: gcsURI},
		},
	}

	op, err := client.LongRunningRecognize(ctx, req)
	if err != nil {
		return err
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	// Print the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Fprintf(out, "\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
			for _, w := range alt.Words {
				fmt.Fprintf(out,
					"Word: \"%v\" (startTime=%3f, endTime=%3f)\n",
					w.Word,
					float64(w.StartTime.Seconds)+float64(w.StartTime.Nanos)*1e-9,
					float64(w.EndTime.Seconds)+float64(w.EndTime.Nanos)*1e-9,
				)
			}
		}
	}
	return nil
}

// [END speech_transcribe_async_time_offsets_gcs]

func syncWords(client *speech.Client, out io.Writer, file string) error {
	ctx := context.Background()

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// Send the contents of the audio file with the encoding and
	// and sample rate information to be transcripted.
	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:              speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz:       16000,
			LanguageCode:          "en-US",
			EnableWordTimeOffsets: true,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	}

	resp, err := client.Recognize(ctx, req)
	if err != nil {
		return err
	}

	// Print the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Fprintf(out, "\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
			for _, w := range alt.Words {
				fmt.Fprintf(out,
					"Word: \"%v\" (startTime=%3f, endTime=%3f)\n",
					w.Word,
					float64(w.StartTime.Seconds)+float64(w.StartTime.Nanos)*1e-9,
					float64(w.EndTime.Seconds)+float64(w.EndTime.Nanos)*1e-9,
				)
			}
		}
	}
	return nil
}
