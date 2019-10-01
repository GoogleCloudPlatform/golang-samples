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

// [START getting_started_auth_setup]

// The authenticating-users program is a sample web server application that
// extracts and verifies user identity data passed to it via Identity-Aware
// Proxy.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/compute/metadata"
	"github.com/dgrijalva/jwt-go"
)

func main() {
	http.HandleFunc("/", index)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// [END getting_started_auth_setup]

// [START getting_started_auth_front_controller]

// index responds to requests with our greeting.
func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	assertion := r.Header.Get("X-Goog-IAP-JWT-Assertion")
	if assertion == "" {
		fmt.Fprintln(w, "Hello None")
		return
	}
	email, _, err := validateAssertion(assertion)
	if err != nil {
		log.Printf("Assertion did not validate: %s", err)
		fmt.Fprintln(w, "Hello None")
		return
	}

	fmt.Fprintf(w, "Hello %s\n", email)
}

// [END getting_started_auth_front_controller]

// [START getting_started_auth_validate]

// validateAssertion validates assertion was signed by Google and returns the
// associated email and userID.
func validateAssertion(assertion string) (email string, userID string, err error) {
	certificates, err := certs()
	if err != nil {
		return "", "", err
	}

	token, err := jwt.Parse(assertion, func(token *jwt.Token) (interface{}, error) {
		keyID := token.Header["kid"].(string)

		_, ok := token.Method.(*jwt.SigningMethodECDSA)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %q", token.Header["alg"])
		}

		cert := certificates[keyID]
		return jwt.ParseECPublicKeyFromPEM([]byte(cert))
	})

	if err != nil {
		return "", "", err
	}

	aud, err := audience()
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("could not extract claims (%T): %+v", token.Claims, token.Claims)
	}

	if claims["aud"].(string) != aud {
		return "", "", fmt.Errorf("mismatched audience. aud field %q does not match %q", claims["aud"], aud)
	}
	return claims["email"].(string), claims["sub"].(string), nil
}

// [END getting_started_auth_validate]

// [START getting_started_auth_audience]

// cachedAudience caches the result of audience.
var cachedAudience string

// audience returns the expected audience value for this service.
func audience() (string, error) {
	if cachedAudience != "" {
		return cachedAudience, nil
	}

	projectNumber, err := metadata.NumericProjectID()
	if err != nil {
		return "", fmt.Errorf("metadata.NumericProjectID: %v", err)
	}

	projectID, err := metadata.ProjectID()
	if err != nil {
		return "", fmt.Errorf("metadata.ProjectID: %v", err)
	}

	cachedAudience = "/projects/" + projectNumber + "/apps/" + projectID

	return cachedAudience, nil
}

// [END getting_started_auth_audience]

// [START getting_started_auth_certs]

// cachedCertificates caches the result of certs.
var cachedCertificates map[string]string

// certs returns IAP's cryptographic public keys.
func certs() (map[string]string, error) {
	if len(cachedCertificates) != 0 { // Already got them previously
		return cachedCertificates, nil
	}

	const url = "https://www.gstatic.com/iap/verify/public_key"
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Get: %v", err)
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&cachedCertificates); err != nil {
		return nil, fmt.Errorf("Decode: %v", err)
	}

	return cachedCertificates, nil
}

// [END getting_started_auth_certs]
