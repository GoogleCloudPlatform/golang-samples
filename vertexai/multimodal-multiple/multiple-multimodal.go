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
package main

// [START aiplatform_gemini_single_turn_multi_image]
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

	if projectID == "" {
		log.Fatal("require environment variable GOOGLE_CLOUD_PROJECT")
	}

	// construct this multimodal prompt:
	// [image of colosseum] city: Rome, Landmark: the Colosseum
	// [image of forbidden city]  city: Beijing, Landmark: the Forbidden City
	// [new image]

	// create prompt image parts
	// colosseum
	colosseum, err := partFromImageURL("https://storage.googleapis.com/cloud-samples-data/vertex-ai/llm/prompts/landmark1.png")
	if err != nil {
		log.Fatalf("unable to read image: %v", err)
	}
	// forbidden city
	forbiddenCity, err := partFromImageURL("https://storage.googleapis.com/cloud-samples-data/vertex-ai/llm/prompts/landmark2.png")
	if err != nil {
		log.Fatalf("unable to read image: %v", err)
	}
	// new image
	newImage, err := partFromImageURL("https://storage.googleapis.com/cloud-samples-data/vertex-ai/llm/prompts/landmark3.png")
	if err != nil {
		log.Fatalf("unable to read image: %v", err)
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
	err = generateMultimodalContent(os.Stdout, prompt, projectID, location, modelName, float32(temperature))
	if err != nil {
		log.Fatalf("unable to generate: %v", err)
	}
}

// generateMultimodalContent provide a generated response using multimodal input
func generateMultimodalContent(w io.Writer, parts []genai.Part, projectID, location, modelName string, temperature float32) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)
	model.Temperature = temperature

	res, err := model.GenerateContent(ctx, parts...)
	if err != nil {
		return fmt.Errorf("unable to generate contents: %v", err)
	}

	fmt.Fprintf(w, "generated response: %s\n", res.Candidates[0].Content.Parts[0])

	return nil
}

// partFromImageURL create a multimodal prompt part from an image URL
func partFromImageURL(image string) (genai.Part, error) {
	var img genai.Blob

	imageURL, err := url.Parse(image)
	if err != nil {
		return img, err
	}
	res, err := http.Get(image)
	if err != nil || res.StatusCode != 200 {
		return img, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return img, fmt.Errorf("unable to read from http: %v", err)
	}

	position := strings.LastIndex(imageURL.Path, ".")
	if position == -1 {
		return img, fmt.Errorf("couldn't find a period to indicate a file extension")
	}
	ext := imageURL.Path[position+1:]

	img = genai.ImageData(ext, data)
	return img, nil
}

// [END aiplatform_gemini_single_turn_multi_image]
