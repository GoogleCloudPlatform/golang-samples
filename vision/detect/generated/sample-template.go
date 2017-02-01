// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

//+build ignore
//go:generate echo foo

// This file is used as the basis for generating detect.go
// To re-generate, run:
//   go generate
// Boilerplate client code is inserted in the sections marked
//   "Boilerplate is inserted by gen.go"

package main

import (
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/vision"
	"golang.org/x/net/context"
)

func init() {
	// Refer to these functions so that goimports is happy before boilerplate is inserted.
	_ = context.Background()
	_ = vision.NewClient
	_ = os.Open
}

// detectFaces gets faces from the Vision API for an image at the given file path.
func detectFaces(w io.Writer, file string) error {
	// Boilerplate is inserted by gen.go
	annotations, err := client.DetectFaces(ctx, image, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No faces found.")
	} else {
		fmt.Fprintln(w, "Faces:")
		for i, annotation := range annotations {
			fmt.Fprintln(w, "  Face ", i)
			fmt.Fprintln(w, "    Anger: ", annotation.Likelihoods.Anger)
			fmt.Fprintln(w, "    Joy: ", annotation.Likelihoods.Joy)
			fmt.Fprintln(w, "    Surprise: ", annotation.Likelihoods.Surprise)
		}
	}

	return nil
}

// detectLabels gets labels from the Vision API for an image at the given file path.
func detectLabels(w io.Writer, file string) error {
	// Boilerplate is inserted by gen.go
	annotations, err := client.DetectLabels(ctx, image, 10)
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
	// Boilerplate is inserted by gen.go
	annotations, err := client.DetectLandmarks(ctx, image, 10)
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
	// Boilerplate is inserted by gen.go
	annotations, err := client.DetectTexts(ctx, image, 10)
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

// detectProperties gets imge properties from the Vision API for an image at the given file path.
func detectProperties(w io.Writer, file string) error {
	// Boilerplate is inserted by gen.go
	props, err := client.DetectImageProps(ctx, image)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Dominant colors:")
	for _, quantized := range props.DominantColors {
		color := quantized.Color
		fmt.Fprintf(w, "%2.1f%% - #%02x%02x%02x\n", quantized.PixelFraction*100, color.R&0xff, color.G&0xff, color.B&0xff)
	}

	return nil
}

// detectSafeSearch gets imge properties from the Vision API for an image at the given file path.
func detectSafeSearch(w io.Writer, file string) error {
	// Boilerplate is inserted by gen.go
	props, err := client.DetectSafeSearch(ctx, image)
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

// detectLogos gets logos from the Vision API for an image at the given file path.
func detectLogos(w io.Writer, file string) error {
	// Boilerplate is inserted by gen.go
	annotations, err := client.DetectLogos(ctx, image, 10)
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
