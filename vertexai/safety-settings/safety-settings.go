// Copyright 2023 Google LLC
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

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/vertexai/genai"
)

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	location := "us-central1"
	modelName := "gemini-pro"
	temperature := 0.4

	prompt := "say something nice to me, but be angry"

	if projectID == "" {
		log.Fatal("require environment variable GOOGLE_CLOUD_PROJECT")
	}

	err := generateContent(os.Stdout, prompt, projectID, location, modelName, float32(temperature))
	if err != nil {
		fmt.Printf("unable to generate: %v\n", err)
	}
}

// generateContent generates text from prompt and configurations provided.
func generateContent(w io.Writer, prompt, projectID, location, modelName string, temperature float32) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return err
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)
	model.Temperature = temperature

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

	res, err := model.GenerateContent(ctx, genai.Text(prompt))
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
