// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_firebase_remote_config]

// Package helloworld provides a set of Cloud Function samples.
package helloworld

import (
	"context"
	"log"
)

// A RemoteConfigEvent is an event triggered by Firebase Remote Config.
type RemoteConfigEvent struct {
	UpdateOrigin string `json:"updateOrigin"`
	UpdateType   string `json:"updateType"`
	UpdateUser   struct {
		Email    string `json:"email"`
		ImageURL string `json:"imageUrl"`
		Name     string `json:"name"`
	} `json:"updateUser"`
	VersionNumber string `json:"versionNumber"`
}

// HelloRemoteConfig handles Firebase Remote Config events.
func HelloRemoteConfig(ctx context.Context, e RemoteConfigEvent) error {
	log.Printf("Update type: %v", e.UpdateType)
	log.Printf("Origin: %v", e.UpdateOrigin)
	log.Printf("Version: %v", e.VersionNumber)
	return nil
}

// [END functions_firebase_remote_config]
