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

// [START eventarc_pubsub_handler]

// Sample pubsub is a Cloud Run service which handles Pub/Sub messages.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

// HelloEventsPubSub receives and processes a Pub/Sub push message.
func HelloEventsPubSub(w http.ResponseWriter, r *http.Request) {
	var e PubSubMessage
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "Bad HTTP Request", http.StatusBadRequest)
		log.Printf("Bad HTTP Request: %v", http.StatusBadRequest)
		return
	}
	name := string(e.Message.Data)
	if name == "" {
		name = "World"
	}
	s := fmt.Sprintf("Hello, %s! ID: %s", name, string(r.Header.Get("Ce-Id")))
	log.Printf(s)
	fmt.Fprintln(w, s)
}

// [END eventarc_pubsub_handler]
// [START eventarc_pubsub_server]

func main() {
	http.HandleFunc("/", HelloEventsPubSub)
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

// [END eventarc_pubsub_server]
