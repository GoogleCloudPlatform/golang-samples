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

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
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

	// Get the full service name from the environment variable
	// set at the time of deployment.
	// Format: "projects/PROJECT_ID/locations/REGION/services/SERVICE_NAME"
	fullServiceName := os.Getenv("FULL_SERVICE_NAME")

	if err := a.getServiceURL(fullServiceName); err != nil {
		return app{}, err
	}

	return a, nil
}

// getServiceURL assigns to internal attribute serviceURL the
// primary URL for a given Cloud Run service.
func (a *app) getServiceURL(fullServiceName string) error {
	ctx := context.Background()

	client, err := run.NewServicesClient(ctx, nil)
	if err != nil {
		return fmt.Errorf("run.NewServicesClient error: %w", err)
	}

	serviceRequest := &runpb.GetServiceRequest{
		Name: fullServiceName,
	}

	service, err := client.GetService(ctx, serviceRequest, nil)
	if err != nil {
		return fmt.Errorf("client.GetService error: %w", err)
	}

	a.serviceURI = service.Uri

	return nil
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

// Parse the authorization header and decode the information beign
// sent by the Bearer Token
func (a *app) receiveAuthorizedRequest(w http.ResponseWriter, r *http.Request) {
	// Allows requests only for the root path ("/") to prevent duplicate calls.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Request method should be GET.
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

	payload, status, err := a.validateToken(token)
	if err != nil {
		http.Error(w, err.Error(), status)
	}

	w.Write(fmt.Appendf(nil, "Hello, %s!\n", payload.Claims["email"]))
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

// [END cloudrun_service_to_service_receive]
