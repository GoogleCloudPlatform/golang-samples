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

// [START video_quickstart]

// Sample video_quickstart uses the Google Cloud Video Intelligence API to label a video.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes"

	video "cloud.google.com/go/videointelligence/apiv1"
	videopb "cloud.google.com/go/videointelligence/apiv1/videointelligencepb"
)

func main() {
	ctx := context.Background()

	// Creates a client.
	client, err := video.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		InputUri: "gs://cloud-samples-data/video/cat.mp4",
		Features: []videopb.Feature{
			videopb.Feature_LABEL_DETECTION,
		},
	})
	if err != nil {
		log.Fatalf("Failed to start annotation job: %v", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		log.Fatalf("Failed to annotate: %v", err)
	}

	// Only one video was processed, so get the first result.
	result := resp.GetAnnotationResults()[0]

	for _, annotation := range result.SegmentLabelAnnotations {
		fmt.Printf("Description: %s\n", annotation.Entity.Description)

		for _, category := range annotation.CategoryEntities {
			fmt.Printf("\tCategory: %s\n", category.Description)
		}

		for _, segment := range annotation.Segments {
			start, _ := ptypes.Duration(segment.Segment.StartTimeOffset)
			end, _ := ptypes.Duration(segment.Segment.EndTimeOffset)
			fmt.Printf("\tSegment: %s to %s\n", start, end)
			fmt.Printf("\tConfidence: %v\n", segment.Confidence)
		}
	}
}

// [END video_quickstart]
