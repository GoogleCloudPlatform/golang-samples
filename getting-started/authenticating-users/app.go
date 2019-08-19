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

	"cloud.google.com/go/compute/metadata"
	"github.com/dgrijalva/jwt-go"
)

var cachedCertificates map[string]string

func main() {
	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func certs() (map[string]string, error) {
	const url = "https://www.gstatic.com/iap/verify/public_key"

	if len(cachedCertificates) != 0 { // Already got them previously
		return cachedCertificates, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Get: %v", err)
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&cachedCertificates); err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}

	return cachedCertificates, nil
}

func audience() (string, error) {
	projectNumber, err := metadata.NumericProjectID()
	if err != nil {
		return "", err
	}

	projectID, err := metadata.ProjectID()
	if err != nil {
		return "", err
	}

	return "/projects/" + projectNumber + "/apps/" + projectID, nil
}

func validateAssertion(assertion string) (email string, userID string, err error) {
	certificates, err := certs()
	if err != nil {
		return "None", "None", err
	}

	token, err := jwt.Parse(assertion, func(token *jwt.Token) (interface{}, error) {
		keyID := token.Header["kid"].(string)

		_, ok := token.Method.(*jwt.SigningMethodECDSA)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		cert := certificates[keyID]
		return jwt.ParseECPublicKeyFromPEM([]byte(cert))
	})

	if err != nil {
		return "None", "None", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "None", "None", fmt.Errorf("could not extract claims")
	}

	aud, err := audience()
	if err != nil {
		return "None", "None", err
	}

	if claims["aud"].(string) != aud {
		return "None", "None", fmt.Errorf("mismatched audience. aud field %s does not match %s", claims["aud"], aud)
	}
	return claims["email"].(string), claims["sub"].(string), nil
}

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	assertion := r.Header.Get("X-Goog-IAP-JWT-Assertion")
	email, _, err := validateAssertion(assertion)
	if err != nil {
		log.Printf("Assertion did not validate: %s", err)
	}

	fmt.Fprintf(w, "Hello %s", email)
}
