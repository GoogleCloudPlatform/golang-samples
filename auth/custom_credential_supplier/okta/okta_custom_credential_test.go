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

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/auth/credentials/externalaccount"
)

// TestOktaSupplier_SubjectToken verifies the fetch logic
// and proper handling of caching.
func TestOktaSupplier_SubjectToken(t *testing.T) {
	// Setup Mock Okta Server
	mockClientID := "client-id"
	mockClientSecret := "client-secret"
	expectedToken := "mock-okta-jwt-token"

	// We count how many times the server is hit to verify caching
	serverHitCount := 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverHitCount++

		// Validate Method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Validate Content Type
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("Expected Content-Type application/x-www-form-urlencoded, got %s", r.Header.Get("Content-Type"))
		}

		// Validate Basic Auth
		authHeader := r.Header.Get("Authorization")
		wantAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(mockClientID+":"+mockClientSecret))
		if authHeader != wantAuth {
			t.Errorf("Invalid Basic Auth header. Got %s, want %s", authHeader, wantAuth)
		}

		// Return Success JSON
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"access_token": expectedToken,
			"expires_in":   3600, // 1 hour
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	// Initialize Supplier with Mock URL
	// Updated struct name to 'oktaSupplier'
	supplier := &oktaSupplier{
		TokenURL:     ts.URL,
		ClientID:     mockClientID,
		ClientSecret: mockClientSecret,
	}
	ctx := context.Background()

	// First Call: Should hit the server
	token, err := supplier.SubjectToken(ctx, &externalaccount.RequestOptions{})
	if err != nil {
		t.Fatalf("First SubjectToken call failed: %v", err)
	}
	if token != expectedToken {
		t.Errorf("Got token %s, want %s", token, expectedToken)
	}
	if serverHitCount != 1 {
		t.Errorf("Expected 1 server hit, got %d", serverHitCount)
	}

	// Second Call: Should use cache (no server hit)
	token, err = supplier.SubjectToken(ctx, &externalaccount.RequestOptions{})
	if err != nil {
		t.Fatalf("Second SubjectToken call failed: %v", err)
	}
	if token != expectedToken {
		t.Errorf("Got token %s, want %s", token, expectedToken)
	}
	if serverHitCount != 1 {
		t.Errorf("Expected server hit count to remain 1 (cached), got %d", serverHitCount)
	}
}

// TestOktaSupplier_ExpiredCache verifies that the supplier
// refetches if the cache is expired.
func TestOktaSupplier_ExpiredCache(t *testing.T) {
	serverHitCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverHitCount++
		w.Header().Set("Content-Type", "application/json")
		// Return a token that expires very soon (e.g., 1 second)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": fmt.Sprintf("token-%d", serverHitCount),
			"expires_in":   30,
		})
	}))
	defer ts.Close()

	supplier := &oktaSupplier{
		TokenURL:     ts.URL,
		ClientID:     "id",
		ClientSecret: "secret",
	}
	ctx := context.Background()

	// First Call
	_, err := supplier.SubjectToken(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if serverHitCount != 1 {
		t.Errorf("Expected 1 hit, got %d", serverHitCount)
	}

	// Second Call immediately after
	// Since expires_in (30s) is < buffer (60s), the supplier considers it expired immediately.
	_, err = supplier.SubjectToken(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if serverHitCount != 2 {
		t.Errorf("Expected 2 hits (cache should be invalid due to short expiry), got %d", serverHitCount)
	}
}

// TestSystem_AuthenticateWithOktaCredentials runs the end-to-end authentication flow
// using values from 'custom-credentials-okta-secrets.json' if present.
func TestSystem_AuthenticateWithOktaCredentials(t *testing.T) {
	const secretsFile = "custom-credentials-okta-secrets.json"

	if _, err := os.Stat(secretsFile); os.IsNotExist(err) {
		t.Skipf("Skipping system test: %s not found", secretsFile)
	}

	// Setup cleanup to restore environment variables after test
	envVars := []string{
		"GCP_WORKLOAD_AUDIENCE",
		"GCS_BUCKET_NAME",
		"GCP_SERVICE_ACCOUNT_IMPERSONATION_URL",
		"OKTA_DOMAIN",
		"OKTA_CLIENT_ID",
		"OKTA_CLIENT_SECRET",
	}
	originalEnv := make(map[string]string)
	for _, k := range envVars {
		originalEnv[k] = os.Getenv(k)
	}
	defer func() {
		for k, v := range originalEnv {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	}()

	loadConfigFromFile()

	audience := os.Getenv("GCP_WORKLOAD_AUDIENCE")
	bucketName := os.Getenv("GCS_BUCKET_NAME")
	oktaDomain := os.Getenv("OKTA_DOMAIN")
	oktaClient := os.Getenv("OKTA_CLIENT_ID")
	oktaSecret := os.Getenv("OKTA_CLIENT_SECRET")

	if audience == "" || bucketName == "" || oktaDomain == "" || oktaClient == "" || oktaSecret == "" {
		t.Skip("Skipping system test: Required configuration missing in secrets file")
	}

	var buf bytes.Buffer
	impersonationURL := os.Getenv("GCP_SERVICE_ACCOUNT_IMPERSONATION_URL")

	err := authenticateWithOktaCredentials(&buf, bucketName, audience, oktaDomain, oktaClient, oktaSecret, impersonationURL)
	if err != nil {
		t.Fatalf("System test failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Success") {
		t.Errorf("Expected output to contain 'Success', got: %s", output)
	}
	if !strings.Contains(output, bucketName) {
		t.Errorf("Expected output to contain bucket name '%s', got: %s", bucketName, output)
	}
}
