// Copyright 2020 Google LLC
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

// [START run_events_pubsub_handler]

// Cloud Run service which handles Pub/Sub messages.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

// HelloEventsPubSub receives and processes a Pub/Sub message via a CloudEvent.
func HelloEventsPubSub(w http.ResponseWriter, r *http.Request) {
	var m PubSubMessage
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var name = string(m.Message.Data)
	if name == "" {
		name = "World"
	}

	s := fmt.Sprintf("Hello, %s! ID: %s", name, string(r.Header.Get("Ce-Id")))
	log.Printf(s)
	fmt.Printf(s)
}

func main() {
	log.Print("run_events_pubsub: starting server...")

	http.HandleFunc("/", HelloEventsPubSub)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("run_events_pubsub: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// [END run_events_pubsub_handler]
