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

// [START run_events_pubsub_server]

// Sample run-pubsub is a Cloud Run service which handles Pub/Sub messages.
package main

import (
	"context"
	"fmt"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func main() {
	ctx := context.Background()
	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Printf("Listening on http://localhost:8080\n")
	log.Fatalf("Failed to start receiver: %s", c.StartReceiver(ctx, HelloPubSub))
}

// PubSubMessage is the payload of the Pub/Sub message.
type PubSubMessage struct {
	Data []byte `json:"data,omitempty"`
	ID   string `json:"id"`
}

// PubSub is the whole payload of a Pub/Sub event.
type PubSub struct {
	Message      PubSubMessage `json:"message"`
	Subscription string        `json:"subscription"`
}

// HelloPubSub receives and processes a Pub/Sub CloudEvent.
func HelloPubSub(ctx context.Context, event cloudevents.Event) {
	// Try to decode the request body into the struct.
	var m PubSub
	err := event.DataAs(&m)
	if err != nil {
		// Error parsing CloudEvent
		log.Fatalf("event.DataAs: could not read CloudEvent: %v", err)
	}
	// Print the data from the Pub/Sub CloudEvent.
	name := string(m.Message.Data)
	if name == "" {
		name = "World"
	}
	s := fmt.Sprintf("Hello, %s! ID: %s", name, event.ID())
	log.Print(s)
}

// [END run_events_pubsub_server]
