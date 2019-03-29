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

// [START videointelligence_text_detection]

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"context"

	"github.com/golang/protobuf/ptypes"

	video "cloud.google.com/go/videointelligence/apiv1"
	videopb "google.golang.org/genproto/googleapis/cloud/videointelligence/v1"
)

// textDetection analyzes a video and extracts the text from the video's audio.
func textDetection(w io.Writer, filename string) error {
	// filename := "resources/googlework_short.mp4"

	ctx := context.Background()

	// Creates a client.
	client, err := video.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		InputContent: fileBytes,
		Features: []videopb.Feature{
			videopb.Feature_TEXT_DETECTION,
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

// [END videointelligence_text_detection]
