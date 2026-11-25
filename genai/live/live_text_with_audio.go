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
	"strings"

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

	modelName := "gemini-3.8-live"

	// The model responds with audio. Request a transcript of it to get text.
	config := &genai.LiveConnectConfig{
		ResponseModalities:       []genai.Modality{genai.ModalityAudio},
		OutputAudioTranscription: &genai.AudioTranscriptionConfig{},
		// Mark the start and end of the audio yourself, because the audio file
		// doesn't end with the silence that automatic activity detection needs.
		RealtimeInputConfig: &genai.RealtimeInputConfig{
			AutomaticActivityDetection: &genai.AutomaticActivityDetection{
				Disabled: true,
			},
		},
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

	// Send the audio as Blob input with explicit activity boundaries
	if err := session.SendRealtimeInput(genai.LiveRealtimeInput{
		ActivityStart: &genai.ActivityStart{},
	}); err != nil {
		return fmt.Errorf("failed to send activity start: %w", err)
	}
	err = session.SendRealtimeInput(genai.LiveRealtimeInput{
		Audio: &genai.Blob{
			Data:     audioBytes,
			MIMEType: "audio/pcm;rate=16000",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send audio input: %w", err)
	}
	if err := session.SendRealtimeInput(genai.LiveRealtimeInput{
		ActivityEnd: &genai.ActivityEnd{},
	}); err != nil {
		return fmt.Errorf("failed to send activity end: %w", err)
	}

	// Stream the response
	var response strings.Builder
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

		if chunk.ServerContent.OutputTranscription != nil && chunk.ServerContent.OutputTranscription.Text != "" {
			response.WriteString(chunk.ServerContent.OutputTranscription.Text)
		}
		if chunk.ServerContent.TurnComplete {
			break
		}
	}

	fmt.Fprintln(w, response.String())

	// Example output:
	// > Answer to this audio url: https://storage.googleapis.com/generativeai-downloads/data/16000.wav
	// Yes, I can hear you. How can I help you today?
	return nil
}

// [END googlegenaisdk_live_txt_with_audio]
