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

// oktaClientCredentialsSupplier implements externalaccount.SubjectTokenProvider
type oktaSupplier struct {
	TokenURL, ClientID, ClientSecret string
	mu                               sync.Mutex
	cachedToken                      string
	expiry                           time.Time
}

// SubjectToken returns a valid Okta access token, refreshing if necessary.
func (s *oktaSupplier) SubjectToken(ctx context.Context, _ *externalaccount.RequestOptions) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cachedToken != "" && time.Now().Add(60*time.Second).Before(s.expiry) {
		return s.cachedToken, nil
	}

	token, expires, err := s.fetchToken(ctx)
	if err != nil {
		return "", err
	}
	s.cachedToken = token
	s.expiry = time.Now().Add(time.Duration(expires) * time.Second)
	return s.cachedToken, nil
}

func (s *oktaSupplier) fetchToken(ctx context.Context) (string, int64, error) {
	v := url.Values{"grant_type": {"client_credentials"}, "scope": {"gcp.test.read"}}
	req, _ := http.NewRequestWithContext(ctx, "POST", s.TokenURL, strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(s.ClientID, s.ClientSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("okta status: %d", resp.StatusCode)
	}

	var res struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", 0, err
	}
	return res.AccessToken, res.ExpiresIn, nil
}

// authenticateWithOktaCredentials authenticates with Google Cloud and verifies bucket access.
func authenticateWithOktaCredentials(w io.Writer, bucket, audience, domain, id, secret, impURL string) error {
	ctx := context.Background()
	tokenURL := fmt.Sprintf("%s/oauth2/default/v1/token", strings.TrimRight(domain, "/"))

	// Initialize credentials with the custom supplier
	creds, err := externalaccount.NewCredentials(&externalaccount.Options{
		Audience:                       audience,
		SubjectTokenType:               "urn:ietf:params:oauth:token-type:jwt",
		ServiceAccountImpersonationURL: impURL,
		SubjectTokenProvider:           &oktaSupplier{TokenURL: tokenURL, ClientID: id, ClientSecret: secret},
		Scopes:                         []string{"https://www.googleapis.com/auth/devstorage.read_only"},
	})
	if err != nil {
		return fmt.Errorf("NewCredentials: %w", err)
	}

	// Get the OAuth2 token
	token, err := creds.Token(ctx)
	if err != nil {
		return fmt.Errorf("Token: %w", err)
	}

	// Verify access via raw HTTP request to GCS
	gcsURL := fmt.Sprintf("https://storage.googleapis.com/storage/v1/b/%s", bucket)
	req, _ := http.NewRequestWithContext(ctx, "GET", gcsURL, nil)
	req.Header.Set("Authorization", "Bearer "+token.Value)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("GCS request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GCS returned %d: %s", resp.StatusCode, string(b))
	}

	fmt.Fprintf(w, "Success! Authenticated and accessed bucket: %s\n", bucket)
	return nil
}

// [END auth_custom_credential_supplier_okta]

func main() {
	loadConfigFromFile()

	audience := os.Getenv("GCP_WORKLOAD_AUDIENCE")
	bucket := os.Getenv("GCS_BUCKET_NAME")
	domain := os.Getenv("OKTA_DOMAIN")
	id := os.Getenv("OKTA_CLIENT_ID")
	secret := os.Getenv("OKTA_CLIENT_SECRET")
	impURL := os.Getenv("GCP_SERVICE_ACCOUNT_IMPERSONATION_URL")

	if audience == "" || bucket == "" || domain == "" || id == "" || secret == "" {
		fmt.Fprintln(os.Stderr, "Missing required configuration. Check custom-credentials-okta-secrets.json or env vars.")
		os.Exit(1)
	}

	if err := authenticateWithOktaCredentials(os.Stdout, bucket, audience, domain, id, secret, impURL); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func loadConfigFromFile() {
	file, err := os.Open("custom-credentials-okta-secrets.json")
	if err != nil {
		return
	}
	defer file.Close()

	var s struct {
		Audience string `json:"gcp_workload_audience"`
		Bucket   string `json:"gcs_bucket_name"`
		ImpURL   string `json:"gcp_service_account_impersonation_url"`
		Domain   string `json:"okta_domain"`
		ID       string `json:"okta_client_id"`
		Secret   string `json:"okta_client_secret"`
	}
	if json.NewDecoder(file).Decode(&s) == nil {
		setEnv("OKTA_DOMAIN", s.Domain)
		setEnv("OKTA_CLIENT_ID", s.ID)
		setEnv("OKTA_CLIENT_SECRET", s.Secret)
		setEnv("GCP_WORKLOAD_AUDIENCE", s.Audience)
		setEnv("GCS_BUCKET_NAME", s.Bucket)
		setEnv("GCP_SERVICE_ACCOUNT_IMPERSONATION_URL", s.ImpURL)
	}
}

func setEnv(k, v string) {
	if v != "" {
		os.Setenv(k, v)
	}
}
