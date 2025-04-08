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

// tokencount shows an example of determining how many tokens correspond to
// a given prompt string
package tokencount

// [START generativeaionvertexai_gemini_token_count_multimodal]
import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"

	"cloud.google.com/go/vertexai/genai"
)

// countTokensMultimodal finds the number of tokens for a multimodal prompt (video+text), and writes to w. Then,
// it calls the model with the multimodal prompt and writes token counts from the response metadata to w.
//
// video is a Google Cloud Storage path starting with "gs://"
func countTokensMultimodal(w io.Writer, projectID, location, modelName string) error {
	// location := "us-central1"
	// modelName := "gemini-1.5-flash-001"
	prompt := "Provide a description of the video."
	video := "gs://cloud-samples-data/generative-ai/video/pixel8.mp4"

	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	part1 := genai.Text(prompt)

	// Given a video file URL, prepare video file as genai.Part
	part2 := genai.FileData{
		MIMEType: mime.TypeByExtension(filepath.Ext(video)),
		FileURI:  video,
	}

	// Finds the total number of tokens for the 2 parts (text, video) of the multimodal prompt,
	// before actually calling the model for inference.
	resp, err := model.CountTokens(ctx, part1, part2)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Number of tokens for the multimodal video prompt: %d\n", resp.TotalTokens)

	res, err := model.GenerateContent(ctx, part1, part2)
	if err != nil {
		return fmt.Errorf("unable to generate contents: %w", err)
	}

	// The token counts are also provided in the model response metadata, after inference.
	fmt.Fprintln(w, "\nModel response")
	md := res.UsageMetadata
	fmt.Fprintf(w, "Prompt Token Count: %d\n", md.PromptTokenCount)
	fmt.Fprintf(w, "Candidates Token Count: %d\n", md.CandidatesTokenCount)
	fmt.Fprintf(w, "Total Token Count: %d\n", md.TotalTokenCount)

	return nil
}

// [END generativeaionvertexai_gemini_token_count_multimodal]
