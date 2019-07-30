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

// The main program is a sample web server application that extracts and
// verifies user identity data passed to it via Identity-Aware Proxy.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

var cachedAudience string
var cachedCertificates = make(map[string]string)

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

func certs() map[string]string {
	const url = "https://www.gstatic.com/iap/verify/public_key"

	if len(cachedCertificates) != 0 {  // Already got them previously
        return cachedCertificates
    }

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch certificates: %s", err)
		return cachedCertificates
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&cachedCertificates)
	if err != nil {
		log.Printf("Error converting from JSON: %s", err)
		return cachedCertificates
	}

	return cachedCertificates
}

func getMetadata(itemName string) string {
	const url = "http://metadata.google.internal/computeMetadata/v1/project/"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url+itemName, nil)
	req.Header.Add("Metadata-Flavor", "Google")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making metadata request: %s", err)
		return "None"
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Error reading metadata: %s", err)
		return "None"
	} else {
		return (string(body))
	}
}

func audience() string {
	if cachedAudience == "" {
		projectNumber := getMetadata("numeric-project-id")
		projectID := getMetadata("project-id")
		cachedAudience = "/projects/" + projectNumber + "/apps/" + projectID
	}

	return cachedAudience
}

func validateAssertion(assertion string) (email string, userid string) {
	certificates := certs()

	token, err := jwt.Parse(assertion, func(token *jwt.Token) (interface{}, error) {
		keyID := token.Header["kid"].(string)

		_, ok := token.Method.(*jwt.SigningMethodECDSA)
		if !ok {
			log.Printf("Wrong signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		cert := certificates[keyID]
		return jwt.ParseECPublicKeyFromPEM([]byte(cert))
	})

	if err != nil {
		log.Printf("Failed to validate assertion: %s", assertion)
		return "None", "None"
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	if claims["aud"].(string) != audience() {
		log.Printf("Token aud %s does not match audience %s", claims["aud"], audience())
		return "None", "None"
	}
	return claims["email"].(string), claims["sub"].(string)
}

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	assertion := r.Header.Get("X-Goog-IAP-JWT-Assertion")
	email, _ := validateAssertion(assertion)
	fmt.Fprint(w, "Hello "+email)
}
