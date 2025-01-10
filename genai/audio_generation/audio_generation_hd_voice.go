// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package audio_generation

// [START genai_audio_generation_hd_voice]
import (
	"context"
	"fmt"
	"io"
	"os"

	genai "google.golang.org/genai"
)

// generateAudioWithHDVoice generates a synthetic text-to-speech audio file with a preset voice.
func generateAudioWithHDVoice(w io.Writer) error {
	modelName := "gemini-2.0-flash-exp"
	client, err := genai.NewClient(context.TODO(), &genai.ClientConfig{})
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"AUDIO"},
		SpeechConfig: &genai.SpeechConfig{
			VoiceConfig: &genai.VoiceConfig{
				PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
					VoiceName: "Kore",
				},
			},
		},
	}
	instance := genai.Text("Say Hello World, I'm Gemini")
	result, err := client.Models.GenerateContent(
		context.TODO(), modelName, instance, config)
	if err != nil {
		return fmt.Errorf("GenerateContent: %w", err)
	}

	part := result.Candidates[0].Content.Parts[0]
	inlineData := part.InlineData
	if inlineData == nil {
		return fmt.Errorf("inlineData is nil")
	}

	audioData := inlineData.Data
	outfile := "audio.wav"
	err = os.WriteFile(outfile, audioData, 0644)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Audio content written to file: %v\n", outfile)
	return nil
}

// [END genai_audio_generation_hd_voice]
