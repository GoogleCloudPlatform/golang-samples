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

// [START googlegenaisdk_live_conversation_audio_with_audio]
import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"google.golang.org/genai"
)

// generateLiveAudioConversation demonstrates two-way audio interaction with a Gemini model using live streaming.
func generateLiveAudioConversation(w io.Writer, audioFilePath string) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{
			APIVersion: "v1beta1",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	modelName := "gemini-live-2.5-flash-preview-native-audio-09-2025"

	// Configure model to receive and respond with audio, including transcriptions.
	config := &genai.LiveConnectConfig{
		ResponseModalities:       []genai.Modality{genai.ModalityAudio},
		InputAudioTranscription:  &genai.AudioTranscriptionConfig{},
		OutputAudioTranscription: &genai.AudioTranscriptionConfig{},
	}

	session, err := client.Live.Connect(ctx, modelName, config)
	if err != nil {
		return fmt.Errorf("failed to connect live: %w", err)
	}
	defer session.Close()

	// Load the audio file
	audioBytes, mimeType, err := loadAudioAsPCMBytes(audioFilePath)
	if err != nil {
		return fmt.Errorf("failed to load audio: %w", err)
	}

	fmt.Fprintf(w, "> Streaming audio from %s to the model\n\n", audioFilePath)

	// Send audio data to the model
	err = session.SendRealtimeInput(genai.LiveSendRealtimeInputParameters{
		Media: &genai.Blob{
			Data:     audioBytes,
			MIMEType: mimeType,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send realtime input: %w", err)
	}

	// Gather audio response frames
	var audioFrames [][]byte

	for {
		chunk, err := session.Receive()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error receiving response: %w", err)
		}

		if chunk.ServerContent != nil {
			if chunk.ServerContent.InputTranscription != nil {
				fmt.Fprintf(w, "Input transcription: %s\n", chunk.ServerContent.InputTranscription.Text)
			}
			if chunk.ServerContent.OutputTranscription != nil {
				fmt.Fprintf(w, "Output transcription: %s\n", chunk.ServerContent.OutputTranscription.Text)
			}
			if chunk.ServerContent.ModelTurn != nil {
				for _, part := range chunk.ServerContent.ModelTurn.Parts {
					if part.InlineData != nil && len(part.InlineData.Data) > 0 {
						audioFrames = append(audioFrames, part.InlineData.Data)
					}
				}
			}
		}
	}

	// Save audio frames to WAV file if available
	if len(audioFrames) > 0 {
		outputFile := "model_response.wav"
		err := saveAudioFramesAsWAV(outputFile, audioFrames, 24000)
		if err != nil {
			return fmt.Errorf("failed to write WAV: %w", err)
		}
		fmt.Fprintf(w, "Model response saved to %s\n", outputFile)
	}

	// Example output:
	// gemini-2.0-flash-live-preview-04-09
	// {'input_transcription': {'text': 'Hello.'}}
	// {'output_transcription': {}}
	// {'output_transcription': {'text': 'Hi'}}
	// {'output_transcription': {'text': ' there. What can I do for you today?'}}
	// {'output_transcription': {'finished': True}}
	// Model response saved to example_model_response.wav
	return nil
}

// loadAudioAsPCMBytes reads a WAV file and returns PCM bytes with a MIME type.
func loadAudioAsPCMBytes(path string) ([]byte, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open WAV file: %w", err)
	}
	defer file.Close()

	wavDecoder := wav.NewDecoder(file)
	if !wavDecoder.IsValidFile() {
		return nil, "", fmt.Errorf("invalid WAV file")
	}
	buf, err := wavDecoder.FullPCMBuffer()
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode WAV: %w", err)
	}

	sampleRate := wavDecoder.SampleRate
	rawInts := buf.Data
	data := make([]byte, len(rawInts)*2) // 16-bit PCM

	for i, sample := range rawInts {
		binary.LittleEndian.PutUint16(data[i*2:], uint16(int16(sample)))
	}

	mimeType := fmt.Sprintf("audio/pcm;rate=%d", sampleRate)
	return data, mimeType, nil
}

// saveAudioFramesAsWAV writes audio frames (PCM bytes) to a WAV file.
func saveAudioFramesAsWAV(filePath string, frames [][]byte, sampleRate int) error {
	audioData := []byte{}
	for _, f := range frames {
		audioData = append(audioData, f...)
	}

	// Create buffer
	intData := audio.IntBuffer{
		Format: &audio.Format{NumChannels: 1, SampleRate: sampleRate},
		Data:   make([]int, len(audioData)/2),
	}

	for i := 0; i < len(audioData); i += 2 {
		intData.Data[i/2] = int(int16(audioData[i]) | int16(audioData[i+1])<<8)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create WAV file: %w", err)
	}
	defer file.Close()

	wavEncoder := wav.NewEncoder(file, sampleRate, 16, 1, 1)
	if err := wavEncoder.Write(&intData); err != nil {
		return fmt.Errorf("failed to write audio data: %w", err)
	}

	if err := wavEncoder.Close(); err != nil {
		return fmt.Errorf("failed to finalize WAV file: %w", err)
	}

	return nil
}

// [END googlegenaisdk_live_conversation_audio_with_audio]
