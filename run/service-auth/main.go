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
	"io"
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
	serviceURI string
}

// newApp returns an app with the serviceURI attribute assigned.
func newApp() (app, error) {
	a := app{}

	if err := a.getServiceURL(); err != nil {
		return app{}, err
	}

	return a, nil
}

// getServiceURL assigns to internal attribute serviceURL the deployed
// service URL.
func (a *app) getServiceURL() error {

	ctx := context.Background()

	// Get the Service Name as found in Cloud Run.
	serviceName := os.Getenv("K_SERVICE")

	// Get the Project ID.
	projectRequest, err := http.NewRequest(http.MethodGet, "http://metadata.google.internal/computeMetadata/v1/project/project-id", nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest error: %w", err)
	}
	projectRequest.Header.Set("Metadata-Flavor", "Google")

	projectResponse, err := http.DefaultClient.Do(projectRequest)
	if err != nil {
		return fmt.Errorf("http.DefaultClient.Do error: %w", err)
	}
	defer projectResponse.Body.Close()

	resBody, err := io.ReadAll(projectResponse.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll error: %w", err)
	}
	projectId := string(resBody)

	// Get the Region.
	regionRequest, err := http.NewRequest(http.MethodGet, "http://metadata.google.internal/computeMetadata/v1/instance/region", nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest error: %w", err)
	}
	regionRequest.Header.Set("Metadata-Flavor", "Google")

	regionResponse, err := http.DefaultClient.Do(regionRequest)
	if err != nil {
		return fmt.Errorf("http.DefaultClient.Do error: %w", err)
	}
	defer regionResponse.Body.Close()

	resBody, err = io.ReadAll(regionResponse.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll error: %w", err)
	}

	splitBody := strings.Split(string(resBody), "/")
	region := splitBody[3]

	// Build fullServiceName.
	fullServiceName := fmt.Sprintf("projects/%s/locations/%s/services/%s", projectId, region, serviceName)

	// Get deployed service URI.
	client, err := run.NewServicesClient(ctx)
	if err != nil {
		return fmt.Errorf("run.NewServicesClient error: %w", err)
	}

	service, err := client.GetService(ctx, &runpb.GetServiceRequest{
		Name: fullServiceName,
	})
	if err != nil {
		return fmt.Errorf("client.GetService error: %w", err)
	}

	// Assign deployed service's URI to internal attribute serviceURL.
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

	// validate token.
	payload, err := validator.Validate(ctx, token, a.serviceURI)
	if err != nil {
		return nil, http.StatusUnauthorized, fmt.Errorf("invalid token: %v", err)
	}

	return payload, http.StatusOK, nil
}

// Parse the authorization header and decode the information beign
// sent by the Bearer Token
func (a *app) receiveAuthorizedRequest(w http.ResponseWriter, r *http.Request) {
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
