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

//+build ignore //# omit

package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	video "cloud.google.com/go/videointelligence/apiv1"
	"github.com/golang/protobuf/ptypes"
	videopb "google.golang.org/genproto/googleapis/cloud/videointelligence/v1"
)

//# if gcs
//# replace __SUFFIX__( URI(
//# end
//# if !gcs
//# replace __SUFFIX__( (
//# end
func boilerplate() { //# omit
	//# def dorequest
	ctx := context.Background()
	client, err := video.NewClient(ctx)
	if err != nil {
		return err
	}
	//# if !gcs

	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	//# end

	op, err := client.AnnotateVideo(ctx, &videopb.AnnotateVideoRequest{
		Features: []videopb.Feature{
			__req.feature__,
		},
		InputContent: fileBytes, //# include if !gcs
		InputUri:     file,      //# include if gcs
	})
	if err != nil {
		return err
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return err
	}
	//# enddef
} //# omit

// [START video_analyze_labels] //# include if !gcs
// [START video_analyze_labels_gcs] //# include if gcs

func label__SUFFIX__(w io.Writer, file string) error {
	//# replace __req.feature__ videopb.Feature_LABEL_DETECTION
	var resp *videopb.AnnotateVideoResponse //# template dorequest

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

// [END video_analyze_labels] //# include if !gcs
// [END video_analyze_labels_gcs] //# include if gcs

// [START video_analyze_shots] //# include if gcs

func shotChange__SUFFIX__(w io.Writer, file string) error {
	//# replace __req.feature__ videopb.Feature_SHOT_CHANGE_DETECTION
	var resp *videopb.AnnotateVideoResponse //# template dorequest

	// A single video was processed. Get the first result.
	result := resp.AnnotationResults[0].ShotAnnotations

	for _, shot := range result {
		start, _ := ptypes.Duration(shot.StartTimeOffset)
		end, _ := ptypes.Duration(shot.EndTimeOffset)

		fmt.Fprintf(w, "Shot: %s to %s\n", start, end)
	}

	return nil
}

// [END video_analyze_shots] //# include if gcs

// [START video_analyze_explicit_content] //# include if gcs

func explicitContent__SUFFIX__(w io.Writer, file string) error {
	//# replace __req.feature__ videopb.Feature_EXPLICIT_CONTENT_DETECTION
	var resp *videopb.AnnotateVideoResponse //# template dorequest

	// A single video was processed. Get the first result.
	result := resp.AnnotationResults[0].ExplicitAnnotation

	for _, frame := range result.Frames {
		offset, _ := ptypes.Duration(frame.TimeOffset)
		fmt.Fprintf(w, "%s - %s\n", offset, frame.PornographyLikelihood.String())
	}

	return nil
}

// [END video_analyze_explicit_content] //# include if gcs

// [START video_analyze_speech_transcription] //# include if !gcs
// [START video_analyze_speech_transcription_gcs] //# include if gcs

func speechTranscription__SUFFIX__(w io.Writer, file string) error {
	ctx := context.Background()
	client, err := video.NewClient(ctx)
	if err != nil {
		return err
	}
	//# if !gcs

	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	//# end

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
		InputContent: fileBytes, //# include if !gcs
		InputUri:     file,      //# include if gcs
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

// [END video_analyze_speech_transcription] //# include if !gcs
// [END video_analyze_speech_transcription_gcs] //# include if gcs
