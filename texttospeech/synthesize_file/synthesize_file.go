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

// The synthesize_file command converts a plain text or SSML file to an audio file.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

// [START tts_synthesize_text_file]

// SynthesizeTextFile synthesizes the text in textFile and saves the output to
// outputFile.
func SynthesizeTextFile(w io.Writer, textFile, outputFile string) error {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	text, err := os.ReadFile(textFile)
	if err != nil {
		return err
	}

	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: string(text)},
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

	err = os.WriteFile(outputFile, resp.AudioContent, 0644)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Audio content written to file: %v\n", outputFile)
	return nil
}

// [END tts_synthesize_text_file]

// [START tts_synthesize_ssml_file]

// SynthesizeSSMLFile synthesizes the SSML contents in ssmlFile and saves the
// output to outputFile.
//
// ssmlFile must be well-formed according to:
//
//	https://www.w3.org/TR/speech-synthesis/
//
// Example: <speak>Hello there.</speak>
func SynthesizeSSMLFile(w io.Writer, ssmlFile, outputFile string) error {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	ssml, err := os.ReadFile(ssmlFile)
	if err != nil {
		return err
	}

	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Ssml{Ssml: string(ssml)},
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

	err = os.WriteFile(outputFile, resp.AudioContent, 0644)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Audio content written to file: %v\n", outputFile)
	return nil
}

// [END tts_synthesize_ssml_file]

func main() {
	ssmlFile := flag.String("ssml", "",
		"The ssml file string from which to synthesize speech.")
	outputFile := flag.String("output-file", "output.mp3",
		"The name of the output file.")
	flag.Parse()

	if *ssmlFile != "" {
		err := SynthesizeSSMLFile(os.Stdout, *ssmlFile, *outputFile)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(`Error: please supply a --text or --ssml content.

Examples:
  go run synthesize_file.go --text ../resources/hello.txt
  go run synthesize_file.go --ssml ../resources/hello.ssml`)
	}
}
