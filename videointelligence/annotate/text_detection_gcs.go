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

// Package annotate contains speech examples.
package annotate

// [START videointelligence_text_detection_gcs]

import (
	"context"
	"fmt"
	"io"

	video "cloud.google.com/go/videointelligence/apiv1"
	"github.com/golang/protobuf/ptypes"
	videopb "google.golang.org/genproto/googleapis/cloud/videointelligence/v1"
)

// textDetectionGCS analyzes a video and extracts the text from the video's audio.
func textDetectionGCS(w io.Writer, gcsURI string) error {
	// gcsURI := "gs://python-docs-samples-tests/video/googlework_short.mp4"

	ctx := context.Background()

	// Creates a client.
	client, err := video.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("video.NewClient: %v", err)
	}

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		InputUri: gcsURI,
		Features: []videopb.Feature{
			videopb.Feature_TEXT_DETECTION,
		},
	})
	if err != nil {
		return fmt.Errorf("AnnotateVideo: %v", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %v", err)
	}

	// Only one video was processed, so get the first result.
	result := resp.GetAnnotationResults()[0]

	for _, annotation := range result.TextAnnotations {
		fmt.Fprintf(w, "Text: %q\n", annotation.GetText())

		// Get the first text segment.
		segment := annotation.GetSegments()[0]
		start, _ := ptypes.Duration(segment.GetSegment().GetStartTimeOffset())
		end, _ := ptypes.Duration(segment.GetSegment().GetEndTimeOffset())
		fmt.Fprintf(w, "\tSegment: %v to %v\n", start, end)

		fmt.Fprintf(w, "\tConfidence: %f\n", segment.GetConfidence())

		// Show the result for the first frame in this segment.
		frame := segment.GetFrames()[0]
		seconds := float32(frame.GetTimeOffset().GetSeconds())
		nanos := float32(frame.GetTimeOffset().GetNanos())
		fmt.Fprintf(w, "\tTime offset of the first frame: %fs\n", seconds+nanos/1e9)

		fmt.Fprintf(w, "\tRotated bounding box vertices:\n")
		for _, vertex := range frame.GetRotatedBoundingBox().GetVertices() {
			fmt.Fprintf(w, "\t\tVertex x=%f, y=%f\n", vertex.GetX(), vertex.GetY())
		}
	}

	return nil
}

// [END videointelligence_text_detection_gcs]
