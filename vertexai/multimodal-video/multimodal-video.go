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

// multimodal video shows generation of content with video and text input
package main

// [START aiplatform_gemini_single_turn_video]
import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
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

	// create video part
	video, err := partFromGCSURI("gs://cloud-samples-data/video/animals.mp4")
	if err != nil {
		log.Fatalf("unable to process media: %v", err)
	}

	// assemble a multimodal (multipart) prompt
	prompt := []genai.Part{
		genai.Text("What is in the video? "),
		video,
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
		return fmt.Errorf("unable to create client: %v", err)
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

// partFromGCSURI create a multimodal prompt part from a Google Cloud Storage object URI
func partFromGCSURI(gcsPath string) (genai.Part, error) {
	gcsPath = strings.Replace(gcsPath, "gs://", "", -1)
	bucket := strings.SplitN(gcsPath, "/", 3)[0]
	object := strings.Join(strings.SplitN(gcsPath, "/", 3)[1:], "/")

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %w", object, err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %w", err)
	}

	position := strings.LastIndex(object, ".")
	if position == -1 {
		return nil, fmt.Errorf("couldn't find a period to indicate a file extension")
	}

	ext := object[position+1:]

	return genai.Blob{MIMEType: "video/" + ext, Data: data}, nil
}

// [END aiplatform_gemini_single_turn_video]
