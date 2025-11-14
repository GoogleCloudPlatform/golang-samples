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
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloud.google.com/go/auth/credentials/externalaccount"
)

// TestOktaClientCredentialsSupplier_SubjectToken verifies the fetch logic
// and proper handling of caching.
func TestOktaClientCredentialsSupplier_SubjectToken(t *testing.T) {
	// 1. Setup Mock Okta Server
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

		// Validate Body Params
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.FormValue("grant_type") != "client_credentials" {
			t.Errorf("Expected grant_type=client_credentials, got %s", r.FormValue("grant_type"))
		}
		if r.FormValue("scope") != "gcp.test.read" {
			t.Errorf("Expected scope=gcp.test.read, got %s", r.FormValue("scope"))
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

	// 2. Initialize Supplier with Mock URL
	supplier := &oktaClientCredentialsSupplier{
		TokenURL:     ts.URL,
		ClientID:     mockClientID,
		ClientSecret: mockClientSecret,
	}
	ctx := context.Background()

	// 3. First Call: Should hit the server
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

	// 4. Second Call: Should use cache (no server hit)
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

// TestOktaClientCredentialsSupplier_ExpiredCache verifies that the supplier
// refetches if the cache is expired.
func TestOktaClientCredentialsSupplier_ExpiredCache(t *testing.T) {
	serverHitCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverHitCount++
		w.Header().Set("Content-Type", "application/json")
		// Return a token that expires very soon (e.g., 1 second)
		// The logic has a 60s buffer, so anything < 60s should be treated as expired immediately for next call
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": fmt.Sprintf("token-%d", serverHitCount),
			"expires_in":   30,
		})
	}))
	defer ts.Close()

	supplier := &oktaClientCredentialsSupplier{
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
