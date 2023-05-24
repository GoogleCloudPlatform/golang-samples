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

package main

// [START cloudrun_service_to_service_receive]

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

const port = 8080

func main() {
	http.HandleFunc("/", receiveAuthorizedGetRequest)
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("Defaulting to port %s", port)
	}
	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// receiveAuthorizedGetRequest takes the "Authorization" header from a
// request, decodes it using the jwt-go library, and returns back the email
// from the header to the caller.
func receiveAuthorizedGetRequest(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		fmt.Fprintf(w, "Hello, anonymous user.\n")
		return
	}

	authValues := strings.SplitN(authHeader, " ", 2)
	if len(authValues) != 2 {
		fmt.Fprintf(w, "Unhandled header format (%v).\n", authHeader)
		return
	}

	authType, creds := authValues[0], authValues[1]

	if authType == "Bearer" {
		ctx := r.Context()
		config := &oauth2.Config{
			ClientID:     "<your-client-id>",
			ClientSecret: "<your-client-secret>",
			Endpoint:     oauth2.Endpoint{
				TokenURL: "https://oauth2.googleapis.com/token",
			},
		}

		token, err := config.TokenSource(ctx, &oauth2.Token{AccessToken: creds}).Token()
		if err != nil {
			http.Error(w, "Unable to parse token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		email := token.Extra("email").(string)
		fmt.Fprintf(w, "Hello, %v!\n", email)
	} else {
		fmt.Fprintf(w, "Unhandled header format (%v).\n", authType)
	}
}

// [END cloudrun_service_to_service_receive]
