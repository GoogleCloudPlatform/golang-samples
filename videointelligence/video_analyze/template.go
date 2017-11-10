// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

//+build ignore //# omit

package main

import (
	"fmt"
	"io"
	"io/ioutil" //# omit if gcs
	"os"        //# omit if gcs

	video "cloud.google.com/go/videointelligence/apiv1"
	videopb "google.golang.org/genproto/googleapis/cloud/videointelligence/v1"

	"github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
)

func boilerplate() error { //# omit
	//# def newclient
	ctx := context.Background()
	client, err := video.NewClient(ctx)
	if err != nil {
		return err
	}

	//# if !gcs
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	//# end
	//# enddef

	var fileBytes []byte              //# omit
	_ = videopb.AnnotateVideoRequest{ //# omit
		//# def input
		InputContent: fileBytes, //# omit if gcs
		InputUri:     file,      //# omit if !gcs
		//# enddef
	} //# omit

	//# def opwait
	if err != nil {
		return err
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return err
	}
	//# enddef
	_ = client //# omit
	_ = resp   //# omit
} //# omit

// label ...
func label(w io.Writer, file string) error {
	var client *video.Client //# replace newclient

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		Features: []videopb.Feature{
			videopb.Feature_LABEL_DETECTION,
		},
		//# replace input
	})
	resp, _ := op.Wait() //# replace opwait

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

// shotChange ...
func shotChange(w io.Writer, file string) error {
	var client *video.Client //# replace newclient

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		Features: []videopb.Feature{
			videopb.Feature_SHOT_CHANGE_DETECTION,
		},
		//# replace input
	})
	resp, _ := op.Wait(ctx) //# replace opwait

	// A single video was processed. Get the first result.
	result := resp.AnnotationResults[0].ShotAnnotations

	for _, shot := range result {
		start, _ := ptypes.Duration(shot.StartTimeOffset)
		end, _ := ptypes.Duration(shot.EndTimeOffset)

		fmt.Fprintf(w, "Shot: %s to %s\n", start, end)
	}

	return nil
}

// explicitContent ...
func explicitContent(w io.Writer, file string) error {
	var client *video.Client //# replace newclient

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		Features: []videopb.Feature{
			videopb.Feature_EXPLICIT_CONTENT_DETECTION,
		},
		//# replace input
	})
	resp, _ := op.Wait(ctx) //# replace opwait

	// A single video was processed. Get the first result.
	result := resp.AnnotationResults[0].ExplicitAnnotation

	for _, frame := range result.Frames {
		offset, _ := ptypes.Duration(frame.TimeOffset)
		fmt.Fprintf(w, "%s - %s\n", offset, frame.PornographyLikelihood.String())
	}

	return nil
}
