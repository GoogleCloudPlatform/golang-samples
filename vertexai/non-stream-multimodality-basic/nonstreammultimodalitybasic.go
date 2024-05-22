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

package nonstreammultimodalitybasic

// [START generativeaionvertexai_stream_multimodality_basic]
import (
	"context"
	"errors"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

// generateContent shows how to send a multi-modal prompt to a model, writing the response to
// the provided io.Writer.
func generateContent(w io.Writer, projectID, modelName string) error {
	// modelName := "gemini-1.5-pro-preview-0409"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, "us-central1")
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	res, err := model.GenerateContent(
		ctx,
		genai.FileData{
			MIMEType: "video/mp4",
			FileURI:  "gs://cloud-samples-data/generative-ai/video/animals.mp4",
		},
		genai.FileData{
			MIMEType: "image/jpeg",
			FileURI:  "gs://cloud-samples-data/generative-ai/image/character.jpg",
		}, genai.Text("Are these video and image correlated?"),
	)
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

// [END generativeaionvertexai_stream_multimodality_basic]
