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

// [START auth_validate_and_decode_bearer_token_on_go]
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

type app struct {
	// serviceURI will be used as the audience value for validating tokens.
	serviceURI string
}

// newApp returns an app with the serviceURI attribute assigned.
func newApp() (app, error) {
	a := app{}

	// Get the service URL from the environment variable
	// set at the time of deployment.
	a.serviceURI = os.Getenv("SERVICE_URL")

	return a, nil
}

// validateToken is used to validate the provided idToken with a known
// Google cert URL.
func (a *app) validateToken(token string) (*idtoken.Payload, int, error) {
	ctx := context.Background()

	// Verify and decode the JWT
	validator, err := idtoken.NewValidator(ctx, option.WithHTTPClient(http.DefaultClient))
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("unable to create Validator")
	}

	// Validate token using serviceURI as audience.
	payload, err := validator.Validate(ctx, token, a.serviceURI)
	if err != nil {
		return nil, http.StatusUnauthorized, fmt.Errorf("invalid token: %v", err)
	}

	return payload, http.StatusOK, nil
}

func main() {
	a, err := newApp()
	if err != nil {
		log.Fatalf("newApp error: %v", err)
	}

	log.Print("starting server...")
	http.HandleFunc("/", a.receiveAuthorizedRequest)

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

// [END auth_validate_and_decode_bearer_token_on_go]
