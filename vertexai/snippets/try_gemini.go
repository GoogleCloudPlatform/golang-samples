// Copyright 2023 Google LLC
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
package snippets

// [START generativeaionvertexai_gemini_image]
// [START aiplatform_gemini_get_started]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

func tryGemini(w io.Writer, projectID string, location string, modelName string) error {
	// location := "us-central1"
	// modelName := "gemini-1.5-flash-001"

	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	gemini := client.GenerativeModel(modelName)

	img := genai.FileData{
		MIMEType: "image/jpeg",
		FileURI:  "gs://generativeai-downloads/images/scones.jpg",
	}
	prompt := genai.Text("What is in this image?")

	resp, err := gemini.GenerateContent(ctx, img, prompt)
	if err != nil {
		return fmt.Errorf("error generating content: %w", err)
	}
	rb, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}
	fmt.Fprintln(w, string(rb))
	return nil
}

// [END aiplatform_gemini_get_started]
// [END generativeaionvertexai_gemini_image]
