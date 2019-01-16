// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_imagemagick_setup]

// Package imagemagick contains an example of using ImageMagick from a Cloud
// Function.
package imagemagick

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// Global API clients used across function invocations.
var (
	storageClient *storage.Client
	visionClient  *vision.ImageAnnotatorClient
)

func init() {
	// Declare a separate err variable to avoid shadowing the client variables.
	var err error

	storageClient, err = storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}

	visionClient, err = vision.NewImageAnnotatorClient(context.Background())
	if err != nil {
		log.Fatalf("vision.NewAnnotatorClient: %v", err)
	}
}

// [END functions_imagemagick_setup]

// [START functions_imagemagick_analyze]

// GCSEvent is the payload of a GCS event. Please refer to the docs for
// additional information regarding GCS events.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

// BlurOffensiveImages blurs offensive images uploaded to GCS.
func BlurOffensiveImages(ctx context.Context, e GCSEvent) error {
	if strings.HasPrefix(e.Name, "blurred-") {
		log.Printf("The image %q is already blurred", e.Name)
		return nil
	}

	img := vision.NewImageFromURI(fmt.Sprintf("gs://%s/%s", e.Bucket, e.Name))

	resp, err := visionClient.DetectSafeSearch(ctx, img, nil)
	if err != nil {
		return fmt.Errorf("AnnotateImage: %v", err)
	}

	if resp.GetAdult() == visionpb.Likelihood_VERY_LIKELY ||
		resp.GetViolence() == visionpb.Likelihood_VERY_LIKELY {
		return blur(ctx, e.Bucket, e.Name)
	}
	log.Printf("The image %q was detected as OK.", e.Name)
	return nil
}

// [END functions_imagemagick_analyze]

// [START functions_imagemagick_blur]

// blur blurs the image stored at gs://bucket/name and stores the result in
// gs://bucket/blurred-name.
func blur(ctx context.Context, bucket, name string) error {
	inputBlob := storageClient.Bucket(bucket).Object(name)
	r, err := inputBlob.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("NewReader: %v", err)
	}

	outputBlob := storageClient.Bucket(bucket).Object("blurred-" + name)
	w := outputBlob.NewWriter(ctx)
	defer w.Close()

	// Use - as input and output to use stdin and stdout.
	cmd := exec.Command("convert", "-", "-blur", "0x8", "-")
	cmd.Stdin = r
	cmd.Stdout = w

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cmd.Run: %v", err)
	}

	log.Printf("Blurred image has been uploaded to %s", outputBlob.ObjectName())

	return nil
}

// [END functions_imagemagick_blur]
