// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_helloworld_storage]

// Package helloworld provides a set of Cloud Function samples.
package helloworld

import (
	"context"
	"log"
)

// GCSEvent is the payload of a GCS event. Please refer to the docs for
// additional information regarding GCS events.
type GCSEvent struct {
	Bucket         string `json:"bucket"`
	Name           string `json:"name"`
	Metageneration string `json:"metageneration"`
	ResourceState  string `json:"resourceState"`
}

// HelloGCS consumes a GCS event.
func HelloGCS(ctx context.Context, e GCSEvent) error {
	if e.ResourceState == "not_exists" {
		log.Printf("File %v deleted.", e.Name)
		return nil
	}
	if e.Metageneration == "1" {
		// The metageneration attribute is updated on metadata changes.
		// The on create value is 1.
		log.Printf("File %v created.", e.Name)
		return nil
	}
	log.Printf("File %v metadata updated.", e.Name)
	return nil
}

// [END functions_helloworld_storage]
