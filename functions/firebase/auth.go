// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_firebase_auth]

// Package firebase contains a Firestore Cloud Function.
package firebase

import (
	"context"
	"log"
	"time"
)

// AuthEvent is the payload of a Firestore Auth event.
type AuthEvent struct {
	Email    string `json:"email"`
	Metadata struct {
		CreatedAt time.Time `json:"createdAt"`
	} `json:"metadata"`
	UID string `json:"uid"`
}

// HelloAuth is triggered by Firestore Auth events.
func HelloAuth(ctx context.Context, e AuthEvent) error {
	log.Printf("Function triggered by creation or deletion of user: %q", e.UID)
	log.Printf("Created at: %v", e.Metadata.CreatedAt)
	if e.Email != "" {
		log.Printf("Email: %q", e.Email)
	}
	return nil
}

// [END functions_firebase_auth]
