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
// multiple-multimodal shows how to generate content from mixed image and text content
package multiplemultimodal

// [START generativeaionvertexai_gemini_single_turn_multi_image]
// [START aiplatform_gemini_single_turn_multi_image]
import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"

	"cloud.google.com/go/vertexai/genai"
)

// generateMultimodalContent shows how to generate a text from a multimodal prompt using the Gemini model,
// writing the response to the provided io.Writer.
func generateMultimodalContent(w io.Writer, projectID, location, modelName string) error {
	// location := "us-central1"
	// modelName := "gemini-1.5-flash-001"
	ctx := context.Background()

	// create prompt image parts
	colosseum := genai.FileData{
		MIMEType: mime.TypeByExtension(filepath.Ext("landmark1.png")),
		FileURI:  "gs://cloud-samples-data/vertex-ai/llm/prompts/landmark1.png",
	}
	forbiddenCity := genai.FileData{
		MIMEType: mime.TypeByExtension(filepath.Ext("landmark2.png")),
		FileURI:  "gs://cloud-samples-data/vertex-ai/llm/prompts/landmark2.png",
	}
	newImage := genai.FileData{
		MIMEType: mime.TypeByExtension(filepath.Ext("landmark3.png")),
		FileURI:  "gs://cloud-samples-data/vertex-ai/llm/prompts/landmark3.png",
	}
	// create a multimodal (multipart) prompt
	prompt := []genai.Part{
		colosseum,
		genai.Text("city: Rome, Landmark: the Colosseum "),
		forbiddenCity,
		genai.Text("city: Beijing, Landmark: the Forbidden City "),
		newImage,
	}

	// generate the response
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	res, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		return fmt.Errorf("unable to generate contents: %w", err)
	}

	fmt.Fprintf(w, "generated response: %s\n", res.Candidates[0].Content.Parts[0])
	return nil
}

// [END aiplatform_gemini_single_turn_multi_image]
// [END generativeaionvertexai_gemini_single_turn_multi_image]
