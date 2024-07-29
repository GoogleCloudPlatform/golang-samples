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

// multimodalall shows an example of understanding a video with audio as multimodal input
package multimodalall

// [START generativeaionvertexai_gemini_all_modalities]
import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"path/filepath"

	"cloud.google.com/go/vertexai/genai"
)

// generateContentFromVideoWithAudio shows how to send a multi-modal prompt to a model, writing the response to
// the provided io.Writer.
func generateContentFromVideoWithAudio(w io.Writer, projectID, location, modelName string) error {
	// location := "us-central1"
	// modelName := "gemini-1.5-flash-001"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	vidPart := genai.FileData{
		MIMEType: mime.TypeByExtension(filepath.Ext("behind_the_scenes_pixel.mp4")),
		FileURI:  "gs://cloud-samples-data/generative-ai/video/behind_the_scenes_pixel.mp4",
	}

	imgPart := genai.FileData{
		MIMEType: mime.TypeByExtension(filepath.Ext("a-man-and-a-dog.png")),
		FileURI:  "gs://cloud-samples-data/generative-ai/image/a-man-and-a-dog.png",
	}

	res, err := model.GenerateContent(ctx, vidPart, imgPart, genai.Text(`
		Watch each frame in the video carefully and answer the questions.
		Only base your answers strictly on what information is available in the video attached.
		Do not make up any information that is not part of the video and do not be too
		verbose, be to the point.

		Questions:
		- When is the moment in the image happening in the video? Provide a timestamp.
		- What is the context of the moment and what does the narrator say about it?
	`))
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

// [END generativeaionvertexai_gemini_all_modalities]
