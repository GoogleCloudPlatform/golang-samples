// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// The synthesize_file command converts a plain text or SSML file to an audio file.
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/net/context"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

// [START SynthesizeTextFile]

// SynthesizeTextFile synthesizes the text in textFile and saves the output to
// outputFile.
func SynthesizeTextFile(w io.Writer, textFile, outputFile string) error {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return err
	}

	text, err := ioutil.ReadFile(textFile)
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

	err = ioutil.WriteFile(outputFile, resp.AudioContent, 0644)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Audio content written to file: %v\n", outputFile)
	return nil
}

// [END SynthesizeTextFile]

// [START SynthesizeSSMLFile]

// SynthesizeSSMLFile synthesizes the SSML contents in ssmlFile and saves the
// output to outputFile.
//
// ssmlFile must be well-formed according to:
//   https://www.w3.org/TR/speech-synthesis/
// Example: <speak>Hello there.</speak>
func SynthesizeSSMLFile(w io.Writer, ssmlFile, outputFile string) error {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return err
	}

	ssml, err := ioutil.ReadFile(ssmlFile)
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

	err = ioutil.WriteFile(outputFile, resp.AudioContent, 0644)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Audio content written to file: %v\n", outputFile)
	return nil
}

// [END SynthesizeSSMLFile]

func main() {
	textFile := flag.String("text", "",
		"The text file from which to synthesize speech.")
	ssmlFile := flag.String("ssml", "",
		"The ssml file string from which to synthesize speech.")
	outputFile := flag.String("output-file", "output.mp3",
		"The name of the output file.")
	flag.Parse()

	if *textFile != "" {
		err := SynthesizeTextFile(os.Stdout, *textFile, *outputFile)
		if err != nil {
			log.Fatal(err)
		}
	} else if *ssmlFile != "" {
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
