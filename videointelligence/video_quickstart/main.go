// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START video_quickstart]

// Sample video_quickstart uses the Google Cloud Video Intelligence API to label a video.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes"

	video "cloud.google.com/go/videointelligence/apiv1"
	videopb "google.golang.org/genproto/googleapis/cloud/videointelligence/v1"
)

func main() {
	ctx := context.Background()

	// Creates a client.
	client, err := video.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		InputUri: "gs://demomaker/cat.mp4",
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
