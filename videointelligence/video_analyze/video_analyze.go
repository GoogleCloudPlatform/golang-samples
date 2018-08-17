// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"io/ioutil"

	video "cloud.google.com/go/videointelligence/apiv1"
	videopb "google.golang.org/genproto/googleapis/cloud/videointelligence/v1"

	"github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
)

// [START video_analyze_labels]

func label(w io.Writer, file string) error {
	ctx := context.Background()
	client, err := video.NewClient(ctx)
	if err != nil {
		return err
	}

	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		Features: []videopb.Feature{
			videopb.Feature_LABEL_DETECTION,
		},
		InputContent: fileBytes,
	})
	if err != nil {
		return err
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	printLabels := func(labels []*videopb.LabelAnnotation) {
		for _, label := range labels {
			fmt.Fprintf(w, "\tDescription: %s\n", label.Entity.Description)
			for _, category := range label.CategoryEntities {
				fmt.Fprintf(w, "\t\tCategory: %s\n", category.Description)
			}
			for _, segment := range label.Segments {
				start, _ := ptypes.Duration(segment.Segment.StartTimeOffset)
				end, _ := ptypes.Duration(segment.Segment.EndTimeOffset)
				fmt.Fprintf(w, "\t\tSegment: %s to %s\n", start, end)
			}
		}
	}

	// A single video was processed. Get the first result.
	result := resp.AnnotationResults[0]

	fmt.Fprintln(w, "SegmentLabelAnnotations:")
	printLabels(result.SegmentLabelAnnotations)
	fmt.Fprintln(w, "ShotLabelAnnotations:")
	printLabels(result.ShotLabelAnnotations)
	fmt.Fprintln(w, "FrameLabelAnnotations:")
	printLabels(result.FrameLabelAnnotations)

	return nil
}

// [END video_analyze_labels]

func shotChange(w io.Writer, file string) error {
	ctx := context.Background()
	client, err := video.NewClient(ctx)
	if err != nil {
		return err
	}

	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		Features: []videopb.Feature{
			videopb.Feature_SHOT_CHANGE_DETECTION,
		},
		InputContent: fileBytes,
	})
	if err != nil {
		return err
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	// A single video was processed. Get the first result.
	result := resp.AnnotationResults[0].ShotAnnotations

	for _, shot := range result {
		start, _ := ptypes.Duration(shot.StartTimeOffset)
		end, _ := ptypes.Duration(shot.EndTimeOffset)

		fmt.Fprintf(w, "Shot: %s to %s\n", start, end)
	}

	return nil
}

func explicitContent(w io.Writer, file string) error {
	ctx := context.Background()
	client, err := video.NewClient(ctx)
	if err != nil {
		return err
	}

	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		Features: []videopb.Feature{
			videopb.Feature_EXPLICIT_CONTENT_DETECTION,
		},
		InputContent: fileBytes,
	})
	if err != nil {
		return err
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	// A single video was processed. Get the first result.
	result := resp.AnnotationResults[0].ExplicitAnnotation

	for _, frame := range result.Frames {
		offset, _ := ptypes.Duration(frame.TimeOffset)
		fmt.Fprintf(w, "%s - %s\n", offset, frame.PornographyLikelihood.String())
	}

	return nil
}
