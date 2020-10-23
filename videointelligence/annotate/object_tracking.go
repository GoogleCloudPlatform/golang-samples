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

// [START video_object_tracking]

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	video "cloud.google.com/go/videointelligence/apiv1"
	"github.com/golang/protobuf/ptypes"
	videopb "google.golang.org/genproto/googleapis/cloud/videointelligence/v1"
)

// objectTracking analyzes a video and extracts entities with their bounding boxes.
func objectTracking(w io.Writer, filename string) error {
	// filename := "../testdata/cat.mp4"

	ctx := context.Background()

	// Creates a client.
	client, err := video.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("video.NewClient: %v", err)
	}

	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		InputContent: fileBytes,
		Features: []videopb.Feature{
			videopb.Feature_OBJECT_TRACKING,
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

	for _, annotation := range result.ObjectAnnotations {
		fmt.Fprintf(w, "Description: %q\n", annotation.Entity.GetDescription())
		if len(annotation.Entity.EntityId) > 0 {
			fmt.Fprintf(w, "\tEntity ID: %q\n", annotation.Entity.GetEntityId())
		}

		segment := annotation.GetSegment()
		start, _ := ptypes.Duration(segment.GetStartTimeOffset())
		end, _ := ptypes.Duration(segment.GetEndTimeOffset())
		fmt.Fprintf(w, "\tSegment: %v to %v\n", start, end)

		fmt.Fprintf(w, "\tConfidence: %f\n", annotation.GetConfidence())

		// Here we print only the bounding box of the first frame in this segment.
		frame := annotation.GetFrames()[0]
		seconds := float32(frame.GetTimeOffset().GetSeconds())
		nanos := float32(frame.GetTimeOffset().GetNanos())
		fmt.Fprintf(w, "\tTime offset of the first frame: %fs\n", seconds+nanos/1e9)

		box := frame.GetNormalizedBoundingBox()
		fmt.Fprintf(w, "\tBounding box position:\n")
		fmt.Fprintf(w, "\t\tleft  : %f\n", box.GetLeft())
		fmt.Fprintf(w, "\t\ttop   : %f\n", box.GetTop())
		fmt.Fprintf(w, "\t\tright : %f\n", box.GetRight())
		fmt.Fprintf(w, "\t\tbottom: %f\n", box.GetBottom())
	}

	return nil
}

// [END video_object_tracking]
