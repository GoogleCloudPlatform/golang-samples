// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package p contains a Cloud Function that processes Firebase
// Authentication events.
package p

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/functions/metadata"
)

// AuthEvent is the payload of a Firebase Auth event.
// Please refer to the docs for additional information
// regarding Firebase Auth events.
type AuthEvent struct {
	Email string `json:"email"`
	UID   string `json:"uid"`
}

// HelloAuth handles changes to Firebase Auth user objects.
func HelloAuth(ctx context.Context, e AuthEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}
	log.Printf("Function triggered by change to: %v", meta.Resource)
	log.Printf("%+v", e)
	return nil
}
