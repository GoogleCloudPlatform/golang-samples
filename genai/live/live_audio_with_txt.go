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

// [START googlegenaisdk_live_audio_with_txt]
import (
	"context"
	"fmt"
	"io"
	"os"

	"google.golang.org/genai"
)

// generateLiveAudioWithText demonstrates using a live Gemini model
// that performs live audio with text and handles responses.
func generateLiveAudioWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.0-flash-live-preview-04-09"

	voiceName := "Aoede"

	config := &genai.LiveConnectConfig{
		ResponseModalities: []genai.Modality{genai.ModalityAudio},
		SpeechConfig: &genai.SpeechConfig{
			VoiceConfig: &genai.VoiceConfig{
				PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
					VoiceName: voiceName,
				},
			},
		},
	}

	// Open a live session
	session, err := client.Live.Connect(ctx, modelName, config)
	if err != nil {
		return fmt.Errorf("failed to connect live: %w", err)
	}
	defer session.Close()

	// Send the text input
	textInput := "Hello? Gemini are you there?"
	fmt.Fprintf(w, "> %s\n\n", textInput)

	err = session.SendClientContent(genai.LiveClientContentInput{
		Turns: []*genai.Content{
			{
				Role: genai.RoleUser,
				Parts: []*genai.Part{
					{Text: textInput},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send content: %w", err)
	}

	// Receive streaming audio chunks
	var audioData []byte
	for {
		chunk, err := session.Receive()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error receiving stream: %w", err)
		}

		if chunk.ServerContent != nil && chunk.ServerContent.ModelTurn != nil {
			for _, part := range chunk.ServerContent.ModelTurn.Parts {
				if part.InlineData != nil {
					audioData = append(audioData, part.InlineData.Data...)
				}
			}
		}
	}

	// Save audio if data received
	if len(audioData) > 0 {
		audioFile := "output.wav"
		if err := os.WriteFile(audioFile, audioData, 0644); err != nil {
			return fmt.Errorf("failed to write WAV file: %w", err)
		}

		fmt.Fprintf(w, "Received audio answer saved to %s\n", audioFile)
	}

	return nil
}

// [END googlegenaisdk_live_audio_with_txt]
