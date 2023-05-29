// Copyright 2023 Google LLC
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

// [START functions_cloudevent_tips_infinite_retries]

// Package tips contains tips for writing Cloud Functions in Go.
package tips

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	functions.CloudEvent("FiniteRetryPubSub", FiniteRetryPubSub)
}

// MessagePublishedData contains the full Pub/Sub message
// See the documentation for more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// FiniteRetryPubSub demonstrates how to avoid inifinite retries.
func FiniteRetryPubSub(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	// Ignore events that are too old.
	expiration := e.Time().Add(10 * time.Second)
	if time.Now().After(expiration) {
		log.Printf("event timeout: halting retries for expired event '%q'", e.ID())
		return nil
	}

	// Add your message processing logic.
	return processTheMessage(msg)
}

// [END functions_cloudevent_tips_infinite_retries]

// processTheMessage is a stand-in for the domain logic of your function.
func processTheMessage(m MessagePublishedData) error {
	// Return an error to retry this function execution.
	return fmt.Errorf("processTheMessage: runtime error")
}
