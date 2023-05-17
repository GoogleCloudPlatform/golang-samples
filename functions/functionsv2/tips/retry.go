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

// [START functions_cloudevent_tips_retry]
package tips

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	functions.CloudEvent("RetryPubSub", RetryPubSub)
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

// RetryPubSub demonstrates how to toggle using retries.
func RetryPubSub(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	name := string(msg.Message.Data)
	if name == "" {
		name = "World"
	}

	// A misconfigured client will stay broken until the function is redeployed.
	client, err := MisconfiguredDataClient()
	if err != nil {
		log.Printf("MisconfiguredDataClient (retry denied):  %v", err)
		// A nil return indicates that the function does not need a retry.
		return nil
	}

	// Runtime error might be resolved with a new attempt.
	if err = FailedWriteOperation(client, name); err != nil {
		log.Printf("FailedWriteOperation (retry expected): %v", err)
		// A non-nil return indicates that a retry is needed.
		return err
	}

	return nil
}

// [END functions_cloudevent_tips_retry]

var misconfigured = true

// MisconfiguredDataClient simulates failure to retrieve a data storage client.
func MisconfiguredDataClient() (interface{}, error) {
	if misconfigured {
		return nil, errors.New("access denied to data service")
	}
	return nil, nil
}

// FailedWriteOperation simulates failure to write data via the provided client
func FailedWriteOperation(client interface{}, data string) error {
	return errors.New("runtime error")
}
