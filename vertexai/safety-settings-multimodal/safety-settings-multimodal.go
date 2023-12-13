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

package main

// [START aiplatform_gemini_safety_settings]
import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	location := "us-central1"
	modelName := "gemini-pro-vision"
	temperature := 0.4

	if projectID == "" {
		log.Fatal("require environment variable GOOGLE_CLOUD_PROJECT")
	}

	cat, _ := partFromImagePath("https://storage.googleapis.com/cloud-samples-data/generative-ai/image/320px-Felis_catus-cat_on_snow.jpg")

	// create a multipart (multimodal) prompt
	prompt := []genai.Part{
		genai.Text("say something nice about this "),
		cat,
	}

	err := generateContent(os.Stdout, prompt, projectID, location, modelName, float32(temperature))
	if err != nil {
		fmt.Printf("unable to generate: %v\n", err)
	}
}

// generateContent generates text from prompt and configurations provided.
func generateContent(w io.Writer, prompt []genai.Part, projectID, location, modelName string, temperature float32) error {
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

	res, err := model.GenerateContent(ctx, prompt...)
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

// partFromImagePath create a multimodal prompt part from an image file path or URL
func partFromImagePath(imagePath string) (genai.Part, error) {
	var img genai.Blob
	var data []byte
	var ext string

	if strings.HasPrefix(imagePath, "https://") {
		imageURL, err := url.Parse(imagePath)
		if err != nil {
			return img, err
		}
		res, err := http.Get(imagePath)
		if err != nil || res.StatusCode != 200 {
			return img, err
		}
		defer res.Body.Close()
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return img, fmt.Errorf("unable to read from http: %v", err)
		}

		position := strings.LastIndex(imageURL.Path, ".")
		if position == -1 {
			return img, fmt.Errorf("couldn't find a period to indicate a file extension")
		}
		ext = imageURL.Path[position+1:]
	} else {
		var err error
		data, err = os.ReadFile(imagePath)
		if err != nil {
			return img, fmt.Errorf("cannot open file: %v", err)
		}
		ext = filepath.Ext(imagePath)
	}

	img = genai.ImageData(ext, data)
	return img, nil
}

// [END aiplatform_gemini_safety_settings]
