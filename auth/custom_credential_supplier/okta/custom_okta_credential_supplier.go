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

// This sample uses cloud.google.com/go/auth, not the older golang.org/x/oauth2/google.
// [START auth_custom_credential_supplier_okta]

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/auth/credentials/externalaccount"
	"cloud.google.com/go/auth/oauth2adapt"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// OktaClientCredentialsSupplier implements externalaccount.SubjectTokenProvider.
// It fetches OIDC tokens from Okta using the Client Credentials grant.
type OktaClientCredentialsSupplier struct {
	TokenURL     string
	ClientID     string
	ClientSecret string

	// Simple in-memory cache for the token.
	mu          sync.Mutex
	cachedToken string
	expiry      time.Time
}

// SubjectToken returns a valid Okta access token, refreshing it if necessary.
func (s *OktaClientCredentialsSupplier) SubjectToken(ctx context.Context, opts *externalaccount.RequestOptions) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if cached token is valid (with a 60-second buffer).
	if s.cachedToken != "" && time.Now().Add(60*time.Second).Before(s.expiry) {
		log.Println("[OktaSupplier] Returning cached token")
		return s.cachedToken, nil
	}

	log.Println("[OktaSupplier] Fetching new token from Okta...")
	// Fetch a new token.
	token, expiresIn, err := s.fetchToken(ctx)
	if err != nil {
		return "", err
	}

	s.cachedToken = token
	s.expiry = time.Now().Add(time.Duration(expiresIn) * time.Second)
	return s.cachedToken, nil
}

func (s *OktaClientCredentialsSupplier) fetchToken(ctx context.Context) (string, int64, error) {
	v := url.Values{}
	v.Set("grant_type", "client_credentials")
	// Adjust scopes as needed for your Okta app configuration.
	// Often 'api://default' or custom scopes are used here.
	// If you get an error from Okta, try removing this line if your app has a default scope.
	v.Set("scope", "gcp.test.read") 

	req, err := http.NewRequestWithContext(ctx, "POST", s.TokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.ClientID, s.ClientSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to fetch Okta token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Okta Error Body: %s", string(body))
		return "", 0, fmt.Errorf("Okta token endpoint returned status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", 0, fmt.Errorf("failed to decode Okta response: %w", err)
	}

	if result.AccessToken == "" {
		return "", 0, fmt.Errorf("Okta response missing access_token")
	}

	return result.AccessToken, result.ExpiresIn, nil
}

func main() {
	ctx := context.Background()

	// 1. Read configuration from environment
	gcpAudience := os.Getenv("GCP_WORKLOAD_AUDIENCE")
	oktaDomain := os.Getenv("OKTA_DOMAIN")
	oktaClientID := os.Getenv("OKTA_CLIENT_ID")
	oktaClientSecret := os.Getenv("OKTA_CLIENT_SECRET")
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	// Optional:
	saImpersonationURL := os.Getenv("GCP_SERVICE_ACCOUNT_IMPERSONATION_URL")

	if gcpAudience == "" || oktaDomain == "" || oktaClientID == "" || oktaClientSecret == "" || projectID == "" {
		log.Fatal("Missing required environment variables. Please set GCP_WORKLOAD_AUDIENCE, OKTA_DOMAIN, OKTA_CLIENT_ID, OKTA_CLIENT_SECRET, and GOOGLE_CLOUD_PROJECT.")
	}

	// 2. Instantiate the custom supplier
	// Note: Adjust the URL path if your Okta org uses a different auth server (e.g., not 'default')
	oktaTokenURL := fmt.Sprintf("%s/oauth2/default/v1/token", strings.TrimRight(oktaDomain, "/"))
	oktaSupplier := &OktaClientCredentialsSupplier{
		TokenURL:     oktaTokenURL,
		ClientID:     oktaClientID,
		ClientSecret: oktaClientSecret,
	}

	// 3. Configure the credentials options
	opts := &externalaccount.Options{
		Audience:                       gcpAudience,
		SubjectTokenType:               "urn:ietf:params:oauth:token-type:jwt",
		SubjectTokenProvider:           oktaSupplier,
		ServiceAccountImpersonationURL: saImpersonationURL,
		Scopes:  []string{"https://www.googleapis.com/auth/cloud-platform"},
	}

	// 4. Create the credentials
	log.Println("Creating Google credentials using Okta supplier...")
	creds, err := externalaccount.NewCredentials(opts)
	if err != nil {
		log.Fatalf("Failed to create credentials: %v", err)
	}

	oauth2Creds := oauth2adapt.Oauth2CredentialsFromAuthCredentials(creds)


	// 5. Use the credentials to make an API call (Listing GCS Buckets)
	log.Println("Authenticating with Google Cloud Storage...")
	storageClient, err := storage.NewClient(ctx, option.WithCredentials(oauth2Creds))
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
	defer storageClient.Close()

	log.Printf("Listing buckets in project %s...\n", projectID)
	it := storageClient.Buckets(ctx, projectID)
	count := 0
	for {
		bkt, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to list buckets (authentication might have failed): %v", err)
		}
		fmt.Printf(" - %s\n", bkt.Name)
		count++
		if count >= 5 {
			fmt.Println("... (stopping after 5 buckets)")
			break
		}
	}
	log.Println("Success!")
}

// [END auth_custom_credential_supplier_okta]
