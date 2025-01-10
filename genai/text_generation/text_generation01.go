// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package text_generation

// [START genai_text_generation01]
import (
	"context"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

// generateContent generates text response for a given input text prompt.
func generateContent(w io.Writer) error {
	modelName := "gemini-2.0-flash-exp"
	client, err := genai.NewClient(context.TODO(), &genai.ClientConfig{})
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"TEXT"},
	}
	textpart := genai.Text(`Write a haiku`)
	result, err := client.Models.GenerateContent(context.TODO(), modelName,
		&genai.ContentParts{textpart}, config)
	if err != nil {
		return fmt.Errorf("GenerateContent: %w", err)
	}

	fmt.Fprintf(w, result.Candidates[0].Content.Parts[0].Text)
	return nil
}

// [END genai_text_generation01]
