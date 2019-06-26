// Copyright 2019 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START run_pubsub_server]

// Sample run-pubsub is a Cloud Run service which handles Pub/Sub messages.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", HelloPubSub)
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// [END run_pubsub_server]

// [START run_pubsub_handler]
// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

// HelloPubSub consumes a Pub/Sub message.
func HelloPubSub(w http.ResponseWriter, r *http.Request) {
	// Parse the Pub/Sub message.
	var m PubSubMessage

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		log.Printf("json.NewDecoder: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := string(m.Message.Data)
	if name == "" {
		name = "World"
	}
	log.Printf("Hello %s!", name)
}

// [END run_pubsub_handler]
