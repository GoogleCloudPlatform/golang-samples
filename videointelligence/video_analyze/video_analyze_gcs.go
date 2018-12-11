// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"io"

	video "cloud.google.com/go/videointelligence/apiv1"
	"github.com/golang/protobuf/ptypes"
	videopb "google.golang.org/genproto/googleapis/cloud/videointelligence/v1"
)

// [START video_analyze_labels_gcs]

func labelURI(w io.Writer, file string) error {
	ctx := context.Background()
	client, err := video.NewClient(ctx)
	if err != nil {
		return err
	}

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		Features: []videopb.Feature{
			videopb.Feature_LABEL_DETECTION,
		},
		InputUri: file,
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

// [END video_analyze_labels_gcs]

// [START video_analyze_shots]

func shotChangeURI(w io.Writer, file string) error {
	ctx := context.Background()
	client, err := video.NewClient(ctx)
	if err != nil {
		return err
	}

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		Features: []videopb.Feature{
			videopb.Feature_SHOT_CHANGE_DETECTION,
		},
		InputUri: file,
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

// [END video_analyze_shots]

// [START video_analyze_explicit_content]

func explicitContentURI(w io.Writer, file string) error {
	ctx := context.Background()
	client, err := video.NewClient(ctx)
	if err != nil {
		return err
	}

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		Features: []videopb.Feature{
			videopb.Feature_EXPLICIT_CONTENT_DETECTION,
		},
		InputUri: file,
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

// [END video_analyze_explicit_content]

// [START video_analyze_speech_transcription_gcs]

func speechTranscriptionURI(w io.Writer, file string) error {
	ctx := context.Background()
	client, err := video.NewClient(ctx)
	if err != nil {
		return err
	}

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		Features: []videopb.Feature{
			videopb.Feature_SPEECH_TRANSCRIPTION,
		},
		VideoContext: &videopb.VideoContext{
			SpeechTranscriptionConfig: &videopb.SpeechTranscriptionConfig{
				LanguageCode:               "en-US",
				EnableAutomaticPunctuation: true,
			},
		},
		InputUri: file,
	})
	if err != nil {
		return err
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	// A single video was processed. Get the first result.
	result := resp.AnnotationResults[0]

	for _, transcription := range result.SpeechTranscriptions {
		// The number of alternatives for each transcription is limited by
		// SpeechTranscriptionConfig.MaxAlternatives.
		// Each alternative is a different possible transcription
		// and has its own confidence score.
		for _, alternative := range transcription.GetAlternatives() {
			fmt.Fprintf(w, "Alternative level information:\n")
			fmt.Fprintf(w, "\tTranscript: %v\n", alternative.GetTranscript())
			fmt.Fprintf(w, "\tConfidence: %v\n", alternative.GetConfidence())

			fmt.Fprintf(w, "Word level information:\n")
			for _, wordInfo := range alternative.GetWords() {
				startTime := wordInfo.GetStartTime()
				endTime := wordInfo.GetEndTime()
				fmt.Fprintf(w, "\t%4.1f - %4.1f: %v (speaker %v)\n",
					float64(startTime.GetSeconds())+float64(startTime.GetNanos())*1e-9, // start as seconds
					float64(endTime.GetSeconds())+float64(endTime.GetNanos())*1e-9,     // end as seconds
					wordInfo.GetWord(),
					wordInfo.GetSpeakerTag())
			}
		}
	}

	return nil
}

// [END video_analyze_speech_transcription_gcs]
