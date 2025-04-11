// Copyright 2020 Google LLC
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

// [START video_detect_logo_gcs]
import (
	"context"
	"fmt"
	"io"
	"time"

	video "cloud.google.com/go/videointelligence/apiv1"
	videopb "cloud.google.com/go/videointelligence/apiv1/videointelligencepb"
	"github.com/golang/protobuf/ptypes"
)

// logoDetectionGCS analyzes a video and extracts logos with their bounding boxes.
func logoDetectionGCS(w io.Writer, gcsURI string) error {
	// gcsURI := "gs://cloud-samples-data/video/googlework_tiny.mp4"

	ctx := context.Background()

	// Creates a client.
	client, err := video.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("video.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*180)
	defer cancel()

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		InputUri: gcsURI,
		Features: []videopb.Feature{
			videopb.Feature_LOGO_RECOGNITION,
		},
	})
	if err != nil {
		return fmt.Errorf("AnnotateVideo: %w", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %w", err)
	}

	// Only one video was processed, so get the first result.
	result := resp.GetAnnotationResults()[0]

	// Annotations for list of logos detected, tracked and recognized in video.
	for _, annotation := range result.LogoRecognitionAnnotations {
		fmt.Fprintf(w, "Description: %q\n", annotation.Entity.GetDescription())
		// Opaque entity ID. Some IDs may be available in Google Knowledge
		// Graph Search API (https://developers.google.com/knowledge-graph/).
		if len(annotation.Entity.EntityId) > 0 {
			fmt.Fprintf(w, "\tEntity ID: %q\n", annotation.Entity.GetEntityId())
		}

		// All logo tracks where the recognized logo appears. Each track
		// corresponds to one logo instance appearing in consecutive frames.
		for _, track := range annotation.Tracks {
			// Video segment of a track.
			segment := track.GetSegment()
			start, _ := ptypes.Duration(segment.GetStartTimeOffset())
			end, _ := ptypes.Duration(segment.GetEndTimeOffset())
			fmt.Fprintf(w, "\tSegment: %v to %v\n", start, end)
			fmt.Fprintf(w, "\tConfidence: %f\n", track.GetConfidence())

			// The object with timestamp and attributes per frame in the track.
			for _, timestampedObject := range track.TimestampedObjects {
				// Normalized Bounding box in a frame, where the object is
				// located.
				box := timestampedObject.GetNormalizedBoundingBox()
				fmt.Fprintf(w, "\tBounding box position:\n")
				fmt.Fprintf(w, "\t\tleft  : %f\n", box.GetLeft())
				fmt.Fprintf(w, "\t\ttop   : %f\n", box.GetTop())
				fmt.Fprintf(w, "\t\tright : %f\n", box.GetRight())
				fmt.Fprintf(w, "\t\tbottom: %f\n", box.GetBottom())

				// Optional. The attributes of the object in the bounding box.
				for _, attribute := range timestampedObject.Attributes {
					fmt.Fprintf(w, "\t\t\tName: %q\n", attribute.GetName())
					fmt.Fprintf(w, "\t\t\tConfidence: %f\n", attribute.GetConfidence())
					fmt.Fprintf(w, "\t\t\tValue: %q\n", attribute.GetValue())
				}
			}

			// Optional. Attributes in the track level.
			for _, trackAttribute := range track.Attributes {
				fmt.Fprintf(w, "\t\tName: %q\n", trackAttribute.GetName())
				fmt.Fprintf(w, "\t\tConfidence: %f\n", trackAttribute.GetConfidence())
				fmt.Fprintf(w, "\t\tValue: %q\n", trackAttribute.GetValue())
			}
		}

		// All video segments where the recognized logo appears. There might be
		// multiple instances of the same logo class appearing in one VideoSegment.
		for _, segment := range annotation.Segments {
			start, _ := ptypes.Duration(segment.GetStartTimeOffset())
			end, _ := ptypes.Duration(segment.GetEndTimeOffset())
			fmt.Fprintf(w, "\tSegment: %v to %v\n", start, end)
		}
	}

	return nil
}

// [END video_detect_logo_gcs]
