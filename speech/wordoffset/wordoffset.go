// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command wordoffset sends audio data to the Google Speech API
// and prints word offset information.
package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
)

const usage = `Usage: wordoffset <audiofile>

Audio file must be a 16-bit signed little-endian encoded
with a sample rate of 16000.

The path to the audio file may be a GCS URI (gs://...).`

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
	defer client.Close()

	if err := sendFunc(client, os.Stdout, os.Args[1]); err != nil {
		log.Fatal(err)
	}
}

// [START speech_transcribe_async_word_time_offsets_gcs]

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

// [END speech_transcribe_async_word_time_offsets_gcs]

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
