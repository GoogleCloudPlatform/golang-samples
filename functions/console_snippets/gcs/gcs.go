// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package p contains a Google Cloud Storage Cloud Function.
package p

import (
	"context"
	"log"
)

// GCSEvent is the payload of a GCS event. Please refer to the docs for
// additional information regarding GCS events.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

// HelloGCS prints a message when a file is changed in a Cloud Storage bucket.
func HelloGCS(ctx context.Context, e GCSEvent) error {
	log.Printf("Processing file: %s", e.Name)
	return nil
}
