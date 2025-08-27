// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package bounding_box shows how to use the GenAI SDK to generate text that adheres to a specific schema.

package bounding_box

// [START googlegenaisdk_boundingbox_with_txt_img]
import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"net/http"

	"google.golang.org/genai"
)

// BoundingBox represents a bounding box with coordinates and label.
type BoundingBox struct {
	Box2D []int  `json:"box_2d"`
	Label string `json:"label"`
}

// plotBoundingBoxes downloads the image and overlays bounding boxes.
func plotBoundingBoxes(imageURI string, boundingBoxes []BoundingBox) error {
	resp, err := http.Get(imageURI)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// Simple red color for bounding boxes
	red := color.RGBA{255, 0, 0, 255}

	for _, bbox := range boundingBoxes {
		// scale normalized coordinates [0–1000] to absolute pixels
		yMin := bbox.Box2D[0] * bounds.Dy() / 1000
		xMin := bbox.Box2D[1] * bounds.Dx() / 1000
		yMax := bbox.Box2D[2] * bounds.Dy() / 1000
		xMax := bbox.Box2D[3] * bounds.Dx() / 1000

		// draw rectangle border
		for x := xMin; x <= xMax; x++ {
			rgba.Set(x, yMin, red)
			rgba.Set(x, yMax, red)
		}
		for y := yMin; y <= yMax; y++ {
			rgba.Set(xMin, y, red)
			rgba.Set(xMax, y, red)
		}
		fmt.Printf("Box: (%d,%d)-(%d,%d) Label: %s\n", xMin, yMin, xMax, yMax, bbox.Label)
	}

	return nil
}

func generateBoundingBoxesWithText(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return err
	}

	imageURI := "https://storage.googleapis.com/generativeai-downloads/images/socks.jpg"

	// Schema definition for []BoundingBox
	schema := &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"box_2d": {
					Type:  genai.TypeArray,
					Items: &genai.Schema{Type: genai.TypeInteger},
				},
				"label": {Type: genai.TypeString},
			},
			Required: []string{"box_2d", "label"},
		},
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{
				Text: "Return bounding boxes as an array with labels. Never return masks. Limit to 25 objects.",
			}},
		},
		Temperature:      float32Ptr(0.5),
		ResponseMIMEType: "application/json",
		ResponseSchema:   schema,
		SafetySettings: []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockThresholdBlockOnlyHigh,
			},
		},
	}

	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{
					FileData: &genai.FileData{
						FileURI:  imageURI,
						MIMEType: "image/jpeg",
					},
				},
				{Text: "Output the positions of the socks with a face. Label according to position in the image."},
			},
		},
	}

	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", contents, config)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, resp)

	// Parse into []BoundingBox
	var boxes []BoundingBox
	if err := json.Unmarshal([]byte(resp.Text()), &boxes); err != nil {
		return fmt.Errorf("failed to parse bounding boxes: %w", err)
	}

	// Example response:
	//	Box: (962,113)-(2158,1631) Label: top left sock with face
	//	Box: (2656,721)-(3953,2976) Label: top right sock with face
	//...

	return plotBoundingBoxes(imageURI, boxes)
}

func float32Ptr(v float32) *float32 { return &v }

// [END googlegenaisdk_boundingbox_with_txt_img]
