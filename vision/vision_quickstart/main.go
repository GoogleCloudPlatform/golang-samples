// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START vision_quickstart]
// Sample vision_quickstart uses the Google Cloud Vision API to label an image.
package main

import (
	"fmt"
	"golang.org/x/net/context"
	"log"
	"os"

	// Imports the Google Cloud Vision API client package
	"cloud.google.com/go/vision"
)

func main() {
	ctx := context.Background()

	// Creates a client
	client, err := vision.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// The name of the image file to annotate
	fileName := "vision/testdata/cat.jpg"

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	defer file.Close()
	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	labels, err := client.DetectLabels(ctx, image, 10)

	fmt.Printf("Labels:\n")
	for _, label := range labels {
		fmt.Println(label.Description)
	}
}

// [END vision_quickstart]
