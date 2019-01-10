// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package snippets contains speech examples.
package snippets

import (
	speech "cloud.google.com/go/speech/apiv1p1beta1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1p1beta1"

	"context"
	"fmt"
	"io"
	"io/ioutil"
)

// [START speech_transcribe_multichannel_beta]

// transcribeMultichannel generates a transcript from a multichannel speech file and tags the speech from each channel.
func transcribeMultichannel(w io.Writer, path string) error {
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("ReadFile: %v", err)
	}

	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:                            speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz:                     44100,
			LanguageCode:                        "en-US",
			AudioChannelCount:                   2,
			EnableSeparateRecognitionPerChannel: true,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	if err != nil {
		return fmt.Errorf("Recognize: %v", err)
	}

	// Print the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Fprintf(w, "Channel %v: %v\n", result.ChannelTag, alt.Transcript)
		}
	}
	return nil
}

// [END speech_transcribe_multichannel_beta]
