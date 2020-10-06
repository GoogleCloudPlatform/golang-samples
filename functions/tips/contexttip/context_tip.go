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

// [START functions_golang_context]
// [START functions_tips_gcp_apis]
// [START functions_pubsub_setup]

// Package contexttip is an example of how to use Pub/Sub and context.Context in
// a Cloud Function.
package contexttip

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
)

// GCP_PROJECT is a user-set environment variable.
var projectID = os.Getenv("GCP_PROJECT")

// client is a global Pub/Sub client, initialized once per instance.
var client *pubsub.Client

func init() {
	// err is pre-declared to avoid shadowing client.
	var err error

	// client is initialized with context.Background() because it should
	// persist between function invocations.
	client, err = pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
}

// [END functions_pubsub_setup]

// [START functions_pubsub_publish]

type publishRequest struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

// PublishMessage publishes a message to Pub/Sub. PublishMessage only works
// with topics that already exist.
func PublishMessage(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the topic name and message.
	p := publishRequest{}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("json.NewDecoder: %v", err)
		http.Error(w, "Error parsing request", http.StatusBadRequest)
		return
	}

	if p.Topic == "" || p.Message == "" {
		s := "missing 'topic' or 'message' parameter"
		log.Println(s)
		http.Error(w, s, http.StatusBadRequest)
		return
	}

	m := &pubsub.Message{
		Data: []byte(p.Message),
	}
	// Publish and Get use r.Context() because they are only needed for this
	// function invocation. If this were a background function, they would use
	// the ctx passed as an argument.
	id, err := client.Topic(p.Topic).Publish(r.Context(), m).Get(r.Context())
	if err != nil {
		log.Printf("topic(%s).Publish.Get: %v", p.Topic, err)
		http.Error(w, "Error publishing message", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Message published: %v", id)
}

// [END functions_pubsub_publish]
// [END functions_tips_gcp_apis]
// [END functions_golang_context]
