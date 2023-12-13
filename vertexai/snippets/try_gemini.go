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

// [START aiplatform_gemini_get_started]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

var projectId = "PROJECT_ID"
var region = "us-central1"

func tryGemini(w io.Writer, projectId string, region string, modelName string) error {

	client, err := genai.NewClient(context.TODO(), projectId, region)
	gemini := client.GenerativeModel("gemini-pro-vision")

	img := genai.FileData{
		MIMEType: "image/jpeg",
		FileURI:  "gs://generativeai-downloads/images/scones.jpg",
	}
	prompt := genai.Text("What is in this image?")
	resp, err := gemini.GenerateContent(context.Background(), img, prompt)
	if err != nil {
		return fmt.Errorf("error generating content: %w", err)
	}
	rb, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintln(w, string(rb))
	return nil
}

// [END aiplatform_gemini_get_started]
