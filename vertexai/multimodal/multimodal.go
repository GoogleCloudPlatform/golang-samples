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

// multimodal shows an example of understanding multimodal input
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	location := "us-central1"
	modelName := "gemini-pro-vision"
	temperature := 0.4

	prompt := "describe what is in this picture"
	image := "https://storage.googleapis.com/cloud-samples-data/generative-ai/image/320px-Felis_catus-cat_on_snow.jpg"

	if projectID == "" {
		log.Fatal("require environment variable GOOGLE_CLOUD_PROJECT")
	}

	err := generateMultimodalContent(os.Stdout, prompt, image, projectID, location, modelName, float32(temperature))
	if err != nil {
		log.Fatalf("unable to generate: %v", err)
	}
}

// generateMultimodalContent returns generated responses based upon the prompt and configurations provided
func generateMultimodalContent(w io.Writer, prompt, image, projectID, location, modelName string, temperature float32) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("unable to create client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)
	model.Temperature = temperature

	// Given an image file URL, prepare image file as genai.Part
	img, err := partFromImageURL(image)
	if err != nil {
		return fmt.Errorf("unable to open image: %v", err)
	}

	res, err := model.GenerateContent(ctx, img, genai.Text(prompt))
	if err != nil {
		return fmt.Errorf("unable to generate contents: %v", err)
	}

	fmt.Fprintf(w, "generated response: %s\n", res.Candidates[0].Content.Parts[0])

	return nil
}

// partFromImageURL create a multimodal prompt part from an image URL
func partFromImageURL(image string) (genai.Part, error) {
	imageURL, err := url.Parse(image)
	if err != nil {
		return nil, err
	}
	res, err := http.Get(image)
	if err != nil || res.StatusCode != 200 {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read from http: %v", err)
	}

	position := strings.LastIndex(imageURL.Path, ".")
	if position == -1 {
		return nil, fmt.Errorf("couldn't find a period to indicate a file extension")
	}
	ext := imageURL.Path[position+1:]

	return genai.ImageData(ext, data), nil
}
