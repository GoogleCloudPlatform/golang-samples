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

package snippets

// [START aiplatform_gemini_translate]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/protobuf/proto"
)

// geminiTranslate shows how translate text, using gemini
func geminiTranslate(w io.Writer, project, location string) error {
	sourceLanguageCode := "en"
	targetLanguageCode := "fr"
	modelName := "gemini-2.0-flash-001"

	ctx := context.Background()
	client, err := genai.NewClient(ctx, project, location)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	model.GenerationConfig = genai.GenerationConfig{
		TopP:            proto.Float32(1),
		TopK:            proto.Int32(32),
		Temperature:     proto.Float32(0.4),
		MaxOutputTokens: proto.Int32(2048),
	}

	// configure the safety settings thresholds
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
	}
	systemInstruction := fmt.Sprintf("Your mission is to translate text from %xs to %s", sourceLanguageCode, targetLanguageCode)

	model.SystemInstruction = &genai.Content{
		Role:  "user",
		Parts: []genai.Part{genai.Text(systemInstruction)},
	}

	prompt := genai.Text("Translate: how are you doing today?")

	resp, err := model.GenerateContent(ctx, prompt)
	if err != nil {
		return fmt.Errorf("error generating content: %w", err)
	}
	// See the JSON response in
	// https://pkg.go.dev/cloud.google.com/go/vertexai/genai#GenerateContentResponse.
	rb, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}
	fmt.Fprintln(w, string(rb))

	return nil
}

// [END aiplatform_gemini_translate]
