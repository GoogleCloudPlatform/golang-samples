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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

var (
	// Messages received by this instance.
	messagesMu sync.Mutex
	messages   [][]byte

	// defaultHTTPClient aliases http.DefaultClient for testing
	defaultHTTPClient = http.DefaultClient
)

const maxMessages = 10

func main() {
	http.HandleFunc("/pubsub/message/list", listHandler)
	http.HandleFunc("/pubsub/message/receive", receiveMessagesHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// pushRequest represents the payload of a Pub/Sub push message.
type pushRequest struct {
	Message      message `json:"message"`
	Subscription string  `json:"subscription"`
}

type message struct {
	Attributes map[string]string `json:"attributes"`
	Data       []byte            `json:"data"`
	ID         string            `json:"messageId"`
}

// [START gae_standard_pubsub_auth_push]
// receiveMessagesHandler validates authentication token and caches the Pub/Sub
// message received.
func receiveMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	pubsubVerificationToken := os.Getenv("PUBSUB_VERIFICATION_TOKEN")
	// Verify that the request originates from the application.
	if token, ok := r.URL.Query()["token"]; !ok || len(token) != 1 || token[0] != pubsubVerificationToken {
		http.Error(w, "Bad token", http.StatusBadRequest)
		return
	}

	// Get the Cloud Pub/Sub-generated JWT in the "Authorization" header.
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(strings.Split(authHeader, " ")) != 2 {
		http.Error(w, "Missing Authorization header", http.StatusBadRequest)
		return
	}
	token := strings.Split(authHeader, " ")[1]
	// Verify and decode the JWT.
	// If you don't need to control the HTTP client used you can use the
	// convenience method idtoken.Validate instead of creating a Validator.
	v, err := idtoken.NewValidator(r.Context(), option.WithHTTPClient(defaultHTTPClient))
	if err != nil {
		http.Error(w, "Unable to create Validator", http.StatusBadRequest)
		return
	}
	// Please change http://example.com to match with the value you are
	// providing while creating the subscription.
	payload, err := v.Validate(r.Context(), token, "http://example.com")
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid Token: %v", err), http.StatusBadRequest)
		return
	}
	if payload.Issuer != "accounts.google.com" && payload.Issuer != "https://accounts.google.com" {
		http.Error(w, "Wrong Issuer", http.StatusBadRequest)

	}

	var pr pushRequest
	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
		return
	}

	messagesMu.Lock()
	defer messagesMu.Unlock()
	// Limit to ten.
	messages = append(messages, pr.Message.Data)
	if len(messages) > maxMessages {
		messages = messages[len(messages)-maxMessages:]
	}

	fmt.Fprint(w, "OK")
}

// [END gae_standard_pubsub_auth_push]

func listHandler(w http.ResponseWriter, r *http.Request) {
	messagesMu.Lock()
	defer messagesMu.Unlock()

	fmt.Fprintln(w, "Messages:")
	for _, v := range messages {
		fmt.Fprintf(w, "Message: %v\n", string(v))
	}
}
