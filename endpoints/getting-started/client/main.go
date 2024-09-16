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
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/auth/credentials"
)

var (
	url                = flag.String("url", "", "The request url. Required.")
	audience           = flag.String("audience", "", "The audience for the JWT. equired")
	serviceAccountFile = flag.String("service-account-file", "", "Path to service account JSON file. Required.")
)

func main() {
	flag.Parse()

	if *audience == "" || *url == "" || *serviceAccountFile == "" {
		fmt.Println("requires: --url, --audience, --service-account-file, --service-account-email")
		os.Exit(1)
	}
	jwt, err := generateJWT(*serviceAccountFile, *audience)
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
func generateJWT(saKeyfile, audience string) (string, error) {
	creds, err := credentials.DetectDefault(&credentials.DetectOptions{
		CredentialsFile:  saKeyfile,
		Audience:         audience,
		UseSelfSignedJWT: true,
	})
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	tok, err := creds.Token(ctx)
	if err != nil {
		return "", err
	}
	return tok.Value, nil
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
