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

// Command client performs authenticated requests against an Endpoints API server.
package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jws"
)

var (
	url                 = flag.String("url", "", "The request url. Required.")
	audience            = flag.String("audience", "", "The audience for the JWT. equired")
	serviceAccountFile  = flag.String("service-account-file", "", "Path to service account JSON file. Required.")
	serviceAccountEmail = flag.String("service-account-email", "", "Path email associated with the service account. Required.")
)

func main() {
	flag.Parse()

	if *audience == "" || *url == "" || *serviceAccountFile == "" || *serviceAccountEmail == "" {
		fmt.Println("requires: --url, --audience, --service-account-file, --service-account-email")
		os.Exit(1)
	}
	jwt, err := generateJWT(*serviceAccountFile, *serviceAccountEmail, *audience, 3600)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := makeJWTRequest(jwt, *url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s Response: %s", *url, resp)
}

// [START endpoints_generate_jwt_sa]

// generateJWT creates a signed JSON Web Token using a Google API Service Account.
func generateJWT(saKeyfile, saEmail, audience string, expiryLength int64) (string, error) {
	now := time.Now().Unix()

	// Build the JWT payload.
	jwt := &jws.ClaimSet{
		Iat: now,
		// expires after 'expiryLength' seconds.
		Exp: now + expiryLength,
		// Iss must match 'issuer' in the security configuration in your
		// swagger spec (e.g. service account email). It can be any string.
		Iss: saEmail,
		// Aud must be either your Endpoints service name, or match the value
		// specified as the 'x-google-audience' in the OpenAPI document.
		Aud: audience,
		// Sub and Email should match the service account's email address.
		Sub:           saEmail,
		PrivateClaims: map[string]interface{}{"email": saEmail},
	}
	jwsHeader := &jws.Header{
		Algorithm: "RS256",
		Typ:       "JWT",
	}

	// Extract the RSA private key from the service account keyfile.
	sa, err := os.ReadFile(saKeyfile)
	if err != nil {
		return "", fmt.Errorf("could not read service account file: %w", err)
	}
	conf, err := google.JWTConfigFromJSON(sa)
	if err != nil {
		return "", fmt.Errorf("could not parse service account JSON: %w", err)
	}
	block, _ := pem.Decode(conf.PrivateKey)
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("private key parse error: %w", err)
	}
	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	// Sign the JWT with the service account's private key.
	if !ok {
		return "", errors.New("private key failed rsa.PrivateKey type assertion")
	}
	return jws.Encode(jwsHeader, jwt, rsaKey)
}

// [END endpoints_generate_jwt_sa]

// [START endpoints_jwt_request]

// makeJWTRequest sends an authorized request to your deployed endpoint.
func makeJWTRequest(signedJWT, url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+signedJWT)
	req.Header.Add("content-type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer response.Body.Close()
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTTP response: %w", err)
	}
	return string(responseData), nil
}

// [END endpoints_jwt_request]
