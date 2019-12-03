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
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

func TestSetEndpoint(t *testing.T) {
	ctx := context.Background()

	const endpoint = "eu-vision.googleapis.com:443"

	client, err := setEndpoint(ctx, endpoint)
	if err != nil {
		t.Fatalf("setEndpoint: %v", err)
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

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Texts:\n")
	for _, text := range texts {
		fmt.Fprintf(&buf, "%v\n", text.GetDescription())
		for _, vertex := range text.GetBoundingPoly().GetVertices() {
			fmt.Fprintf(&buf, "  bounding vertex: %v, %v\n", vertex.GetX(), vertex.GetY())
		}
	}

	if got, want := buf.String(), "System"; !strings.Contains(got, want) {
		t.Errorf("setEndpoint got:\n----\n%s----\nWant to contain:\n----\n%s\n----", got, want)
	}
	if got, want := buf.String(), "bounding vertex:"; !strings.Contains(got, want) {
		t.Errorf("setEndpoint got:\n----\n%s----\nWant to contain:\n----\n%s\n----", got, want)
	}
}
