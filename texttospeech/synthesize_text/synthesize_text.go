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

// The synthesize_text command converts plain text or SSML content to an audio file.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

// [START tts_synthesize_text]

// SynthesizeText synthesizes plain text and saves the output to outputFile.
func SynthesizeText(w io.Writer, text, outputFile string) error {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return err
	}

	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		// Note: the voice can also be specified by name.
		// Names of voices can be retrieved with client.ListVoices().
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_FEMALE,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputFile, resp.AudioContent, 0644)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Audio content written to file: %v\n", outputFile)
	return nil
}

// [END tts_synthesize_text]

// [START tts_synthesize_ssml]

// SynthesizeSSML synthesizes ssml and saves the output to outputFile.
//
// ssml must be well-formed according to:
//   https://www.w3.org/TR/speech-synthesis/
// Example: <speak>Hello there.</speak>
func SynthesizeSSML(w io.Writer, ssml, outputFile string) error {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return err
	}

	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Ssml{Ssml: ssml},
		},
		// Note: the voice can also be specified by name.
		// Names of voices can be retrieved with client.ListVoices().
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_FEMALE,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputFile, resp.AudioContent, 0644)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Audio content written to file: %v\n", outputFile)
	return nil
}

// [END tts_synthesize_ssml]

func main() {
	text := flag.String("text", "",
		"The text from which to synthesize speech.")
	ssml := flag.String("ssml", "",
		"The ssml string from which to synthesize speech.")
	outputFile := flag.String("output-file", "output.txt",
		"The name of the output file.")
	flag.Parse()

	if *text != "" {
		err := SynthesizeText(os.Stdout, *text, *outputFile)
		if err != nil {
			log.Fatal(err)
		}
	} else if *ssml != "" {
		err := SynthesizeSSML(os.Stdout, *ssml, *outputFile)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(`Error: please supply a --text or --ssml content.

Examples:
  go run synthesize_text.go --text "hello"
  go run synthesize_text.go --ssml "<speak>Hello there.</speak>"`)
	}
}
