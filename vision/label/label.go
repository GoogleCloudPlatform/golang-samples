// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command label uses the Vision API's label detection capabilities to find a label based on an image's content.
package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/vision/v1"
)

// findLabels gets labels from the Vision API for an image at the given file path.
func findLabels(file string) ([]string, error) {
	ctx := context.Background()

	// Authenticate to generate a vision service.
	client, err := google.DefaultClient(ctx, vision.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	service, err := vision.New(client)
	if err != nil {
		return nil, err
	}

	// Read the image.
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Construct a label request, encoding the image in base64.
	req := &vision.AnnotateImageRequest{
		Image: &vision.Image{
			Content: base64.StdEncoding.EncodeToString(b),
		},
		Features: []*vision.Feature{{Type: "LABEL_DETECTION"}},
	}
	batch := &vision.BatchAnnotateImagesRequest{
		Requests: []*vision.AnnotateImageRequest{req},
	}
	res, err := service.Images.Annotate(batch).Do()
	if err != nil {
		return nil, err
	}

	labels := make([]string, 0)
	for _, annotation := range res.Responses[0].LabelAnnotations {
		labels = append(labels, annotation.Description)
	}
	return labels, nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-image>\n", filepath.Base(os.Args[0]))
	}
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	labels, err := findLabels(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if len(labels) == 0 {
		fmt.Println("No labels found.")
	} else {
		fmt.Println("Found labels:")
		for _, label := range labels {
			fmt.Println(label)
		}
	}
}
