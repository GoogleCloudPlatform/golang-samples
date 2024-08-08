// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// multimodalaudio shows an example of understanding multimodal input including audio
package multimodalaudio

// [START generativeaionvertexai_gemini_audio_summarization]
// [START generativeaionvertexai_gemini_audio_transcription]
import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"path/filepath"

	"cloud.google.com/go/vertexai/genai"
)

// [END generativeaionvertexai_gemini_audio_transcription]
// [END generativeaionvertexai_gemini_audio_summarization]

// [START generativeaionvertexai_gemini_audio_summarization]
// summarizeAudio shows how to send an audio asset and a text question to a model, writing the response to the
// provided io.Writer.
func summarizeAudio(w io.Writer, projectID, location, modelName string) error {
	// location := "us-central1"
	// modelName := "gemini-1.5-flash-001"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)
	model.SetTemperature(0.4)

	// Given an audio file URL, prepare audio file as genai.Part
	part := genai.FileData{
		MIMEType: mime.TypeByExtension(filepath.Ext("pixel.mp3")),
		FileURI:  "gs://cloud-samples-data/generative-ai/audio/pixel.mp3",
	}

	res, err := model.GenerateContent(ctx, part, genai.Text(`
		Please provide a summary for the audio.
		Provide chapter titles with timestamps, be concise and short, no need to provide chapter summaries.
		Do not make up any information that is not part of the audio and do not be verbose.
	`,
	))
	if err != nil {
		return fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 ||
		len(res.Candidates[0].Content.Parts) == 0 {
		return errors.New("empty response from model")
	}

	fmt.Fprintf(w, "generated summary:\n%s\n", res.Candidates[0].Content.Parts[0])
	return nil
}

// [END generativeaionvertexai_gemini_audio_summarization]

// [START generativeaionvertexai_gemini_audio_transcription]
// transcribeAudio generates a response into w
func transcribeAudio(w io.Writer, projectID, location, modelName string) error {
	// location := "us-central1"
	// modelName := "gemini-1.5-flash-001"

	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	// Optional: set an explicit temperature
	model.SetTemperature(0.4)

	// Given an audio file URL, prepare audio file as genai.Part
	img := genai.FileData{
		MIMEType: mime.TypeByExtension(filepath.Ext("pixel.mp3")),
		FileURI:  "gs://cloud-samples-data/generative-ai/audio/pixel.mp3",
	}

	res, err := model.GenerateContent(ctx, img, genai.Text(`
			Can you transcribe this interview, in the format of timecode, speaker, caption.
			Use speaker A, speaker B, etc. to identify speakers.
	`))
	if err != nil {
		return fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 ||
		len(res.Candidates[0].Content.Parts) == 0 {
		return errors.New("empty response from model")
	}

	fmt.Fprintf(w, "generated transcript:\n%s\n", res.Candidates[0].Content.Parts[0])
	return nil
}

// [END generativeaionvertexai_gemini_audio_transcription]
