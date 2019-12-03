// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

// [START vision_set_endpoint]
import (
	"context"
	"fmt"
	"io"

	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/option"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// setEndpoint changes your endpoint.
func setEndpoint(w io.Writer, endpoint string) error {
	// endpoint := "eu-vision.googleapis.com:443"

	// Create a client with a custom endpoint.
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("NewImageAnnotatorClient: %v", err)
	}
	defer client.Close()

	// Use the client with custom endpoint to detect texts on an image.
	image := &visionpb.Image{
		Source: &visionpb.ImageSource{
			GcsImageUri: "gs://cloud-samples-data/vision/text/screen.jpg",
		},
	}
	maxResults := 3
	texts, err := client.DetectTexts(ctx, image, nil, maxResults)
	if err != nil {
		return fmt.Errorf("DetectTexts: %v", err)
	}

	fmt.Fprintf(w, "Texts:\n")
	for _, text := range texts {
		fmt.Fprintf(w, "%v\n", text.GetDescription())
		for _, vertex := range text.GetBoundingPoly().GetVertices() {
			fmt.Fprintf(w, "  bounding vertex: %v, %v\n", vertex.GetX(), vertex.GetY())
		}
	}

	return nil
}

// [END vision_set_endpoint]
