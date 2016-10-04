// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command label uses the Vision API's label detection capabilities to find a label based on an image's content.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	// [START imports]
	"cloud.google.com/go/vision"
	"golang.org/x/net/context"
	// [END imports]
)

// findLabels gets labels from the Vision API for an image at the given file path.
func findLabels(file string) ([]string, error) {
	// [START init]
	ctx := context.Background()

	// Create the client.
	client, err := vision.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	// [END init]

	// [START request]
	// Open the file.
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return nil, err
	}

	// Perform the request.
	annotations, err := client.DetectLabels(ctx, image, 10)
	if err != nil {
		return nil, err
	}
	// [END request]
	// [START transform]
	var labels []string
	for _, annotation := range annotations {
		labels = append(labels, annotation.Description)
	}
	return labels, nil
	// [END transform]
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
