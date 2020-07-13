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
	"net/http"
	"os"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

var (
	handler = http.DefaultServeMux
)

func main() {
	ctx := context.Background()
	// Create a new HTTP client for CloudEvents
	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatal(err)
	}
	handleFn, err := cloudevents.NewHTTPReceiveHandler(ctx, p, HelloPubSub)
	if err != nil {
		log.Fatal(err)
	}
	handler.Handle("/", handleFn)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// [END run_events_pubsub_server]

// [START run_events_pubsub_handler]

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
func HelloPubSub(ctx context.Context, event cloudevents.Event) (string, error) {
	// Try to decode the request body into the struct.
	var m PubSub
	err := event.DataAs(&m)
	if err != nil {
		// Error parsing CloudEvent
		return "", fmt.Errorf("event.DataAs: could not read CloudEvent: %v", err)
	}
	// Print and return the data from the Pub/Sub CloudEvent.
	s := fmt.Sprintf("Hello, %s! ID: %s", string(m.Message.Data), event.ID())
	log.Printf(s)
	return s, nil
}

// [END run_events_pubsub_handler]
