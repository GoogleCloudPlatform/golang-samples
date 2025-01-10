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

package image_generation

// [START genai_image_generation_interleaved_image_input]
import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	genai "google.golang.org/genai"
)

// generateContent generates image and text outputs for a given input image and text prompt.
func generateContent(w io.Writer) error {
	modelName := "gemini-2.0-flash-exp"
	client, err := genai.NewClient(context.TODO(), &genai.ClientConfig{})
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}

	resp, err := http.Get("https://storage.googleapis.com/cloud-samples-data/generative-ai/image/sofa.jpeg")
	if err != nil {
		return fmt.Errorf("error fetching image: %w", err)
	}
	defer resp.Body.Close()
	imagebytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %w", err)
	}
	imagepart := genai.InlineData{
		Data:     imagebytes,
		MIMEType: "image/png",
	}
	textpart := genai.Text("What other color sofas would work in my space? Can you update the image?")

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"TEXT", "IMAGE"},
	}
	iter := client.Models.GenerateContentStream(context.TODO(), modelName,
		&genai.ContentParts{textpart, imagepart}, config)

	imagesSeen := 0
	for r, err := range iter {
		if err != nil {
			return fmt.Errorf("GenerateContentStream: %w", err)
		}
		for _, p := range r.Candidates[0].Content.Parts {
			if p.Text != "" {
				fmt.Fprintf(w, "Text response: %s", p.Text)
			}
			if p.InlineData != nil {
				filename := fmt.Sprintf("image-%d.png", imagesSeen)
				err := os.WriteFile(filename, p.InlineData.Data, 0644)
				if err != nil {
					return fmt.Errorf("failed to write file %s: %w", filename, err)
				}
				fmt.Fprintf(w, "Image response : %s\n", filename)
				imagesSeen++
			}
		}
	}
	return nil
}

// [END genai_image_generation_interleaved_image_input]
