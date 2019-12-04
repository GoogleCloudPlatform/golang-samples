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

import (
	"context"
	"strings"
	"testing"

	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/option"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

func TestSetEndpoint(t *testing.T) {
	const endpoint = "eu-vision.googleapis.com:443"

	// Run the code sample to check for errors.
	err := setEndpoint(endpoint)
	if err != nil {
		t.Fatalf("setEndpoint: %v", err)
	}

	// Since we're not returning the client from the code sample, we create an equivalent client here.
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		t.Fatalf("NewImageAnnotatorClient: %v", err)
	}
	defer client.Close()

	image := &visionpb.Image{
		Source: &visionpb.ImageSource{
			GcsImageUri: "gs://cloud-samples-data/vision/text/screen.jpg",
		},
	}
	texts, err := client.DetectTexts(ctx, image, nil, 1)
	if err != nil {
		t.Fatalf("DetectTexts: %v", err)
	}

	text := texts[0]
	if got, want := text.GetDescription(), "System"; !strings.Contains(got, want) {
		t.Errorf("text.GetDescription() got:\n----\n%s----\nWant to contain:\n----\n%s\n----", got, want)
	}
	if len(text.GetBoundingPoly().GetVertices()) == 0 {
		t.Errorf("text.GetBoundingPoly().getVertices() must have at least one vertex")
	}
}
