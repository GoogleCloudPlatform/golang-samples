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

// [START googlegenaisdk_live_txt_with_audio]
import (
	"context"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/genai"
)

// generateLiveTextWithAudio demonstrates sending audio to a live session and
// receiving text output. It sends the audio as a Blob inside a genai.LiveRealtimeInput.
func generateLiveTextWithAudio(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-2.0-flash-live-preview-04-09"

	config := &genai.LiveConnectConfig{
		ResponseModalities: []genai.Modality{genai.ModalityText},
	}

	session, err := client.Live.Connect(ctx, modelName, config)
	if err != nil {
		return fmt.Errorf("failed to connect live: %w", err)
	}
	defer session.Close()

	audioURL := "https://storage.googleapis.com/generativeai-downloads/data/16000.wav"
	// Download audio
	resp, err := http.Get(audioURL)
	if err != nil {
		return fmt.Errorf("failed to download audio: %w", err)
	}
	defer resp.Body.Close()

	audioBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read audio: %w", err)
	}

	fmt.Fprintf(w, "> Answer to this audio url: %s\n\n", audioURL)

	// Send the audio as Blob media input
	err = session.SendRealtimeInput(genai.LiveRealtimeInput{
		Media: &genai.Blob{
			Data:     audioBytes,
			MIMEType: "audio/pcm;rate=16000",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send audio input: %w", err)
	}

	// Stream the response
	var response string
	for {
		chunk, err := session.Receive()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error receiving response: %w", err)
		}

		if chunk.ServerContent == nil {
			continue
		}

		// Handle model turn responses
		if chunk.ServerContent.ModelTurn != nil {
			for _, part := range chunk.ServerContent.ModelTurn.Parts {
				if part != nil && part.Text != "" {
					response += part.Text
				}
			}
		}
	}

	fmt.Fprintln(w, response)

	// Example output:
	// > Answer to this audio url: https://storage.googleapis.com/generativeai-downloads/data/16000.wav
	// Yes, I can hear you. How can I help you today?
	return nil
}

// [END googlegenaisdk_live_txt_with_audio]
