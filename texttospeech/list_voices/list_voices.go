// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// The `list_voices` command lists the available Text-to-Speech voices.
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/net/context"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

// [START ListVoices]

// ListVoices lists the available text to speech voices.
func ListVoices(w io.Writer) error {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return err
	}

	// Performs the list voices request.
	resp, err := client.ListVoices(ctx, &texttospeechpb.ListVoicesRequest{})
	if err != nil {
		return err
	}

	for _, voice := range resp.Voices {
		// Display the voice's name. Example: tpc-vocoded
		fmt.Fprintf(w, "Name: %v\n", voice.Name)

		// Display the supported language codes for this voice. Example: "en-US"
		for languageCode := range voice.LanguageCodes {
			fmt.Fprintf(w, "  Supported language: %v\n", languageCode)
		}

		// Display the SSML Voice Gender.
		fmt.Fprintf(w, "  SSML Voice Gender: %v\n", voice.SsmlGender.String())

		// Display the natural sample rate hertz for this voice. Example: 24000
		fmt.Fprintf(w, "  Natural Sample Rate Hertz: %v\n",
			voice.NaturalSampleRateHertz)
	}

	return nil
}

// [END ListVoices]

func main() {
	err := ListVoices(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
