// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START functions_firebase_remote_config]

// Package helloworld provides a set of Cloud Functions samples.
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
