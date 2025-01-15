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

// multimodalvideo shows an example of understanding multimodal input including video
package multimodalvideo

// [START generativeaionvertexai_gemini_video]
// [START aiplatform_gemini_single_turn_video]
import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"path/filepath"

	"cloud.google.com/go/vertexai/genai"
)

// generateMultimodalContent generates a response into w, based upon the prompt and video.
func generateMultimodalContent(w io.Writer, projectID, location, modelName string) error {
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

	// Given a video file URL, prepare video file as genai.Part
	part := genai.FileData{
		MIMEType: mime.TypeByExtension(filepath.Ext("animals.mp4")),
		FileURI:  "gs://cloud-samples-data/video/animals.mp4",
	}

	res, err := model.GenerateContent(ctx, part, genai.Text("What is in this video?"))
	if err != nil {
		return fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 ||
		len(res.Candidates[0].Content.Parts) == 0 {
		return errors.New("empty response from model")
	}

	fmt.Fprintf(w, "generated response: %s\n", res.Candidates[0].Content.Parts[0])
	return nil
}

// [END aiplatform_gemini_single_turn_video]
// [END generativeaionvertexai_gemini_video]
