// Copyright 2025 Google LLC
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
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

func validateToken(token string) (*idtoken.Payload, int, error) {
	ctx := context.Background()

	// Verify and decode the JWT
	validator, err := idtoken.NewValidator(ctx, option.WithHTTPClient(http.DefaultClient))
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("unable to create Validator")
	}

	payload, err := validator.Validate(ctx, token, "")
	if err != nil {
		return nil, http.StatusUnauthorized, fmt.Errorf("invalid token: %v", err)
	}

	return payload, http.StatusOK, nil
}

// Parse the authorization header and decode the information beign
// sent by the Bearer Token
func receiveAuthorizedRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Attempt to retrieve and validate the Authorization header.
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.Write([]byte("Hello, anonymous user\n"))
		return
	}

	if len(strings.Split(authHeader, " ")) != 2 {
		http.Error(w, "Malformed Authorization header", http.StatusBadRequest)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := strings.Split(authHeader, " ")[1]

	payload, status, err := validateToken(token)
	if err != nil {
		http.Error(w, err.Error(), status)
	}

	w.Write(fmt.Appendf(nil, "Hello, %s!\n", payload.Claims["email"]))
}

func main() {
	log.Print("starting server...")
	http.HandleFunc("/", receiveAuthorizedRequest)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// [END cloudrun_service_to_service_receive]
