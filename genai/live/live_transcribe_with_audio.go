// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package live shows how to use the GenAI SDK to generate text with live resources.
package live

// [START googlegenaisdk_live_transcribe_with_audio]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateLiveTranscribeWithAudio demonstrates using a live Gemini model
// that performs live transcribe with audio and handles responses.
func generateLiveTranscribeWithAudio(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-live-2.5-flash-preview-native-audio"

	config := &genai.LiveConnectConfig{
		ResponseModalities:       []genai.Modality{genai.ModalityAudio},
		InputAudioTranscription:  &genai.AudioTranscriptionConfig{},
		OutputAudioTranscription: &genai.AudioTranscriptionConfig{},
	}

	session, err := client.Live.Connect(ctx, modelName, config)
	if err != nil {
		return fmt.Errorf("failed to connect live session: %w", err)
	}
	defer session.Close()

	inputText := "Hello? Gemini are you there?"
	fmt.Fprintf(w, "> %s\n", inputText)

	err = session.SendClientContent(genai.LiveClientContentInput{
		Turns: []*genai.Content{
			{
				Role: genai.RoleUser,
				Parts: []*genai.Part{
					{Text: inputText},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send client content: %w", err)
	}

	var response string

	for {
		message, err := session.Receive()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error receiving stream: %w", err)
		}

		if message.ServerContent != nil {
			if message.ServerContent.ModelTurn != nil {
				fmt.Fprintf(w, "Model turn: %v\n", message.ServerContent.ModelTurn)
			}

			// Input transcription from audio
			if message.ServerContent.InputTranscription != nil {
				if message.ServerContent.InputTranscription.Text != "" {
					fmt.Fprintf(w, "Input transcript: %s\n",
						message.ServerContent.InputTranscription.Text)
				}
			}

			// Output transcription (model generated)
			if message.ServerContent.OutputTranscription != nil {
				if message.ServerContent.OutputTranscription.Text != "" {
					response += message.ServerContent.OutputTranscription.Text
				}
			}
		}
	}
	// Example output:
	//  >  Hello? Gemini are you there?
	//  Yes, I'm here. What would you like to talk about?
	fmt.Fprintln(w, response)
	return nil
}

// [END googlegenaisdk_live_transcribe_with_audio]
