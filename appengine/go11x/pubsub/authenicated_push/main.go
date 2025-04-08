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
	"google.golang.org/api/pubsub/v1"
)

type app struct {
	// pubsubVerificationToken is a shared secret between the the publisher of
	// the message and this application.
	pubsubVerificationToken string

	// Messages received by this instance.
	messagesMu sync.Mutex
	messages   []string

	// defaultHTTPClient aliases http.DefaultClient for testing
	defaultHTTPClient *http.Client
}

const maxMessages = 10

func main() {
	app := &app{
		defaultHTTPClient:       http.DefaultClient,
		pubsubVerificationToken: os.Getenv("PUBSUB_VERIFICATION_TOKEN"),
	}
	http.HandleFunc("/pubsub/message/list", app.listHandler)
	http.HandleFunc("/pubsub/message/receive", app.receiveMessagesHandler)

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
	Message      pubsub.PubsubMessage `json:"message"`
	Subscription string               `json:"subscription"`
}

// [START gae_standard_pubsub_auth_push]
// receiveMessagesHandler validates authentication token and caches the Pub/Sub
// message received.
func (a *app) receiveMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Verify that the request originates from the application.
	// a.pubsubVerificationToken = os.Getenv("PUBSUB_VERIFICATION_TOKEN")
	if token, ok := r.URL.Query()["token"]; !ok || len(token) != 1 || token[0] != a.pubsubVerificationToken {
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
	v, err := idtoken.NewValidator(r.Context(), option.WithHTTPClient(a.defaultHTTPClient))
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
		return
	}

	// IMPORTANT: you should validate claim details not covered by signature
	// and audience verification above, including:
	//   - Ensure that `payload.Claims["email"]` is equal to the expected service
	//     account set up in the push subscription settings.
	//   - Ensure that `payload.Claims["email_verified"]` is set to true.
	if payload.Claims["email"] != "test-service-account-email@example.com" || payload.Claims["email_verified"] != true {
		http.Error(w, "Unexpected email identity", http.StatusBadRequest)
		return
	}

	var pr pushRequest
	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
		return
	}

	a.messagesMu.Lock()
	defer a.messagesMu.Unlock()
	// Limit to ten.
	a.messages = append(a.messages, pr.Message.Data)
	if len(a.messages) > maxMessages {
		a.messages = a.messages[len(a.messages)-maxMessages:]
	}

	fmt.Fprint(w, "OK")
}

// [END gae_standard_pubsub_auth_push]

func (a *app) listHandler(w http.ResponseWriter, r *http.Request) {
	a.messagesMu.Lock()
	defer a.messagesMu.Unlock()

	fmt.Fprintln(w, "Messages:")
	for _, v := range a.messages {
		fmt.Fprintf(w, "Message: %v\n", v)
	}
}
