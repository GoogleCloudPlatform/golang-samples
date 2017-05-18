// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

//+build ignore
//go:generate echo foo

// This file is used as the basis for generating detect.go
// To re-generate, run:
//   go generate
// Boilerplate client code is inserted in the sections marked
// `	var client *vision.Client // Boilerplate is inserted by gen.go`
package main

import (
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/vision/apiv1"
	"golang.org/x/net/context"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

func init() {
	// Refer to these functions so that goimports is happy before boilerplate is inserted.
	_ = context.Background()
	_ = vision.ImageAnnotatorClient{}
	_ = os.Open
	_ = pb.AnnotateImageRequest{}
}

// detectFaces gets faces from the Vision API for an image at the given file path.
func detectFaces(w io.Writer, file string) error {
	var client *vision.ImageAnnotatorClient // Boilerplate is inserted by gen.go
	annotations, err := client.DetectFaces(ctx, image, nil, 10)
	if err != nil {
		return err
	}
	if len(annotations) == 0 {
		fmt.Fprintln(w, "No faces found.")
	} else {
		fmt.Fprintln(w, "Faces:")
		for i, annotation := range annotations {
			fmt.Fprintln(w, "  Face", i)
			fmt.Fprintln(w, "    Anger:", annotation.AngerLikelihood)
			fmt.Fprintln(w, "    Joy:", annotation.JoyLikelihood)
			fmt.Fprintln(w, "    Surprise:", annotation.SurpriseLikelihood)
		}
	}
	return nil
}

// detectLabels gets labels from the Vision API for an image at the given file path.
func detectLabels(w io.Writer, file string) error {
	var client *vision.ImageAnnotatorClient // Boilerplate is inserted by gen.go
	annotations, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No labels found.")
	} else {
		fmt.Fprintln(w, "Labels:")
		for _, annotation := range annotations {
			fmt.Fprintln(w, annotation.Description)
		}
	}

	return nil
}

// detectLandmarks gets landmarks from the Vision API for an image at the given file path.
func detectLandmarks(w io.Writer, file string) error {
	var client *vision.ImageAnnotatorClient // Boilerplate is inserted by gen.go
	annotations, err := client.DetectLandmarks(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No landmarks found.")
	} else {
		fmt.Fprintln(w, "Landmarks:")
		for _, annotation := range annotations {
			fmt.Fprintln(w, annotation.Description)
		}
	}

	return nil
}

// detectText gets text from the Vision API for an image at the given file path.
func detectText(w io.Writer, file string) error {
	var client *vision.ImageAnnotatorClient // Boilerplate is inserted by gen.go
	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No text found.")
	} else {
		fmt.Fprintln(w, "Text:")
		for _, annotation := range annotations {
			fmt.Fprintf(w, "%q\n", annotation.Description)
		}
	}

	return nil
}

// detectDocumentText gets the full document text from the Vision API for an image at the given file path.
func detectDocumentText(w io.Writer, file string) error {
	var client *vision.ImageAnnotatorClient // Boilerplate is inserted by gen.go
	annotation, err := client.DetectDocumentText(ctx, image, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Text:")
	fmt.Fprintf(w, "%q\n", annotation.Text)

	return nil
}

// detectProperties gets image properties from the Vision API for an image at the given file path.
func detectProperties(w io.Writer, file string) error {
	var client *vision.ImageAnnotatorClient // Boilerplate is inserted by gen.go
	props, err := client.DetectImageProperties(ctx, image, nil, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Dominant colors:")
	for _, quantized := range props.DominantColors.Colors {
		color := quantized.Color
		r := int(color.Red) & 0xff
		g := int(color.Green) & 0xff
		b := int(color.Blue) & 0xff
		fmt.Fprintf(w, "%2.1f%% - #%02x%02x%02x\n", quantized.PixelFraction*100, r, g, b)
	}

	return nil
}

// detectCropHints gets suggested croppings the Vision API for an image at the given file path.
func detectCropHints(w io.Writer, file string) error {
	var client *vision.ImageAnnotatorClient // Boilerplate is inserted by gen.go
	res, err := client.CropHints(ctx, image, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Crop hints:")
	for _, hint := range res.CropHints {
		fmt.Fprintf(w, "%v\n", hint.BoundingPoly)
	}

	return nil
}

// detectSafeSearch gets image properties from the Vision API for an image at the given file path.
func detectSafeSearch(w io.Writer, file string) error {
	var client *vision.ImageAnnotatorClient // Boilerplate is inserted by gen.go
	props, err := client.DetectSafeSearch(ctx, image, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Safe Search properties:")
	fmt.Fprintln(w, "Adult:", props.Adult)
	fmt.Fprintln(w, "Medical:", props.Medical)
	fmt.Fprintln(w, "Spoofed:", props.Spoof)
	fmt.Fprintln(w, "Violence:", props.Violence)

	return nil
}

// detectWeb gets image properties from the Vision API for an image at the given file path.
func detectWeb(w io.Writer, file string) error {
	var client *vision.ImageAnnotatorClient // Boilerplate is inserted by gen.go
	web, err := client.DetectWeb(ctx, image, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Web properties:")
	if len(web.FullMatchingImages) != 0 {
		fmt.Fprintln(w, "\tFull image matches:")
		for _, full := range web.FullMatchingImages {
			fmt.Fprintf(w, "\t\t%s\n", full.Url)
		}
	}
	if len(web.PagesWithMatchingImages) != 0 {
		fmt.Fprintln(w, "\tPages with this image:")
		for _, page := range web.PagesWithMatchingImages {
			fmt.Fprintf(w, "\t\t%s\n", page.Url)
		}
	}
	if len(web.WebEntities) != 0 {
		fmt.Fprintln(w, "\tEntities:")
		for _, entity := range web.WebEntities {
			fmt.Fprintf(w, "\t\t%-12s %s\n", entity.EntityId, entity.Description)
		}
	}

	return nil
}

// detectLogos gets logos from the Vision API for an image at the given file path.
func detectLogos(w io.Writer, file string) error {
	var client *vision.ImageAnnotatorClient // Boilerplate is inserted by gen.go
	annotations, err := client.DetectLogos(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No logos found.")
	} else {
		fmt.Fprintln(w, "Logos:")
		for _, annotation := range annotations {
			fmt.Fprintln(w, annotation.Description)
		}
	}

	return nil
}
