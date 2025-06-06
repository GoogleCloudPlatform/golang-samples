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
// safety-settings shows how to adjust safety settings for a generative model

package safetysettings

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

// generateContent generates text from prompt and configurations provided.
func generateContent(w io.Writer, projectID, location, modelName string) error {
	// location := "us-central1"
	// model := "gemini-2.0-flash-001"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return err
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)
	model.SetTemperature(0.8)

	// configure the safety settings thresholds
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockLowAndAbove,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockLowAndAbove,
		},
	}

	res, err := model.GenerateContent(ctx, genai.Text("Hello, say something mean to me."))
	if err != nil {
		return fmt.Errorf("unable to generate content: %v", err)
	}
	fmt.Fprintf(w, "generate-content response: %v\n", res.Candidates[0].Content.Parts[0])

	fmt.Fprintf(w, "safety ratings:\n")
	for _, r := range res.Candidates[0].SafetyRatings {
		fmt.Fprintf(w, "\t%+v\n", r)
	}

	return nil
}
