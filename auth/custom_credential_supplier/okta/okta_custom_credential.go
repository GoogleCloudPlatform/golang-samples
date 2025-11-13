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

// [START auth_custom_credential_supplier_okta]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/auth/credentials/externalaccount"
)

// oktaClientCredentialsSupplier implements externalaccount.SubjectTokenProvider.
// It fetches OIDC tokens from Okta using the Client Credentials grant.
type oktaClientCredentialsSupplier struct {
	TokenURL     string
	ClientID     string
	ClientSecret string

	// Simple in-memory cache for the token.
	mu          sync.Mutex
	cachedToken string
	expiry      time.Time
}

// SubjectToken returns a valid Okta access token, refreshing it if necessary.
func (s *oktaClientCredentialsSupplier) SubjectToken(ctx context.Context, opts *externalaccount.RequestOptions) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if cached token is valid (with a 60-second buffer).
	if s.cachedToken != "" && time.Now().Add(60*time.Second).Before(s.expiry) {
		return s.cachedToken, nil
	}

	// Fetch a new token.
	token, expiresIn, err := s.fetchToken(ctx)
	if err != nil {
		return "", err
	}

	s.cachedToken = token
	s.expiry = time.Now().Add(time.Duration(expiresIn) * time.Second)
	return s.cachedToken, nil
}

// fetchToken performs the HTTP request to Okta.
func (s *oktaClientCredentialsSupplier) fetchToken(ctx context.Context) (string, int64, error) {
	v := url.Values{}
	v.Set("grant_type", "client_credentials")
	// The scope here is specific to the Okta application configuration.
	v.Set("scope", "gcp.test.read")

	req, err := http.NewRequestWithContext(ctx, "POST", s.TokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return "", 0, fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(s.ClientID, s.ClientSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to fetch Okta token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", 0, fmt.Errorf("okta token endpoint returned status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", 0, fmt.Errorf("failed to decode Okta response: %w", err)
	}

	if result.AccessToken == "" {
		return "", 0, fmt.Errorf("okta response missing access_token")
	}

	return result.AccessToken, result.ExpiresIn, nil
}

// authenticateWithOktaCredentials demonstrates how to use a custom Okta credential supplier
// to authenticate with Google Cloud and verify access to a specific bucket.
func authenticateWithOktaCredentials(w io.Writer, bucketName, audience, domain, clientID, clientSecret, impersonationURL string) error {
	// bucketName := "sample-bucket"
	// audience := "//iam.googleapis.com/projects/sample-project/locations/global/workloadIdentityPools/sample-pool/providers/sample-provider"
	// domain := "https://sample.okta.com"
	// clientID := "pqr123"
	// clientSecret := "00124cas62huads68755"
	// [Optional] impersonationURL := "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/myserviceaccount@iam.gserviceaccount.com:generateAccessToken"

	ctx := context.Background()

	// 1. Instantiate the custom supplier
	// Note: Adjust the URL path if your Okta org uses a different auth server (e.g., not 'default')
	oktaTokenURL := fmt.Sprintf("%s/oauth2/default/v1/token", strings.TrimRight(domain, "/"))

	supplier := &oktaClientCredentialsSupplier{
		TokenURL:     oktaTokenURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	// 2. Configure the credentials options
	opts := &externalaccount.Options{
		Audience:                       audience,
		SubjectTokenType:               "urn:ietf:params:oauth:token-type:jwt",
		SubjectTokenProvider:           supplier,
		ServiceAccountImpersonationURL: impersonationURL,
		Scopes:                         []string{"https://www.googleapis.com/auth/devstorage.read_write"},
	}

	// 3. Create the credentials object
	creds, err := externalaccount.NewCredentials(opts)
	if err != nil {
		return fmt.Errorf("externalaccount.NewCredentials: %w", err)
	}

	// 4. Authenticate and Check Bucket Existence via Raw HTTP
	bucketURL := fmt.Sprintf("https://storage.googleapis.com/storage/v1/b/%s", bucketName)
	fmt.Fprintf(w, "Request URL: %s\n", bucketURL)
	fmt.Fprintln(w, "Attempting to make authenticated request to Google Cloud Storage...")

	// Retrieve a valid token
	token, err := creds.Token(ctx)
	if err != nil {
		return fmt.Errorf("creds.Token: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", bucketURL, nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.Value)

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()

	// Check for failure
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("bucket '%s' does not exist (404)", bucketName)
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse and display success output
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("json.Decode: %w", err)
	}

	fmt.Fprintln(w, "\n--- SUCCESS! ---")
	fmt.Fprintln(w, "Successfully authenticated and retrieved bucket data:")
	prettyJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintln(w, string(prettyJSON))

	return nil
}

// [END auth_custom_credential_supplier_okta]

// main is provided for local execution and testing purposes.
func main() {
	gcpAudience := os.Getenv("GCP_WORKLOAD_AUDIENCE")
	oktaDomain := os.Getenv("OKTA_DOMAIN")
	oktaClientID := os.Getenv("OKTA_CLIENT_ID")
	oktaClientSecret := os.Getenv("OKTA_CLIENT_SECRET")
	gcsBucketName := os.Getenv("GCS_BUCKET_NAME")
	// Optional
	saImpersonationURL := os.Getenv("GCP_SERVICE_ACCOUNT_IMPERSONATION_URL")

	if gcpAudience == "" || oktaDomain == "" || oktaClientID == "" || oktaClientSecret == "" || gcsBucketName == "" {
		fmt.Fprintln(os.Stderr, "Missing required environment variables: GCP_WORKLOAD_AUDIENCE, OKTA_DOMAIN, OKTA_CLIENT_ID, OKTA_CLIENT_SECRET, GCS_BUCKET_NAME")
		os.Exit(1)
	}

	if err := authenticateWithOktaCredentials(os.Stdout, gcsBucketName, gcpAudience, oktaDomain, oktaClientID, oktaClientSecret, saImpersonationURL); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
