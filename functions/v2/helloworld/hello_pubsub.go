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

// [START functions_helloworld_pubsub_v2]

// Package helloworld provides a set of Cloud Functions samples.
package helloworld

import (
	"context"
	"log"
)

// A CloudEvent containing the Pub/Sub message.
// See the documentation for more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type CloudEventMessage struct {
	Data []byte `json:"message"`
}

// HelloPubSub consumes a Pub/Sub message.
func HelloPubSub(cem CloudEventMessage) error {
	name := string(cem.message.data) // Automatically decoded from base64.
	if name == "" {
		name = "World"
	}
	log.Printf("Hello, %s!", name)
	return nil
}

// [END functions_helloworld_pubsub_v2]
