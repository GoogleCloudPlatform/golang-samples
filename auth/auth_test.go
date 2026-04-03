// Copyright 2022 Google LLC
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

package snippets

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

func TestAuthSnippets(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	audience := "https://example.com"

	want := "Listed all storage buckets."

	buf := &bytes.Buffer{}

	if err := authenticateExplicitWithAdc(buf); err != nil {
		t.Fatalf("authenticateExplicitWithAdc got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("authenticateExplicitWithAdc got %q, want %q", got, want)
	}

	buf.Reset()

	if err := authenticateImplicitWithAdc(buf, tc.ProjectID); err != nil {
		t.Fatalf("authenticateImplicitWithAdc got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("authenticateImplicitWithAdc got %q, want %q", got, want)
	}

	want = "Generated ID token."
	buf.Reset()

	if err := getIdTokenFromMetadataServer(buf, audience); err != nil {
		t.Fatalf("getIdTokenFromMetadataServer got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("getIdTokenFromMetadataServer got %q, want %q", got, want)
	}

	buf.Reset()

	if err := getIdTokenFromServiceAccount(buf, os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"), audience); err != nil {
		t.Fatalf("getIdTokenFromServiceAccount got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("getIdTokenFromServiceAccount got %q, want %q", got, want)
	}

	buf.Reset()
	want = "Generated OAuth2 token"
	impersonatedServiceAccount := fmt.Sprintf("auth-samples-testing@%s.iam.gserviceaccount.com", tc.ProjectID)
	scope := "https://www.googleapis.com/auth/cloud-platform"

	if err := getAccessTokenFromImpersonatedCredentials(buf, impersonatedServiceAccount, scope); err != nil {
		t.Fatalf("getAccessTokenFromImpersonatedCredentials got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("getAccessTokenFromImpersonatedCredentials got %q, want %q", got, want)
	}

	buf.Reset()
	want = "ID token verified."

	credentials, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		t.Fatalf("failed to generate default credentials: %v", err)
	}

	ts, err := idtoken.NewTokenSource(ctx, audience, option.WithCredentials(credentials))
	if err != nil {
		t.Fatalf("failed to create NewTokenSource: %v", err)
	}

	token, err := ts.Token()
	if err != nil {
		t.Fatalf("failed to get ID token: %v", err)
	}

	if err := verifyGoogleIdToken(buf, token.AccessToken, audience); err != nil {
		t.Fatalf("verifyGoogleIdToken got err: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("verifyGoogleIdToken got %q, want %q", got, want)
	}
}

func TestAuthenticateWithAPIKey(t *testing.T) {
	apiKey := os.Getenv("GOLANG_SAMPLES_API_KEY")
	buf := &bytes.Buffer{}
	if err := authenticateWithAPIKey(buf, apiKey); err != nil {
		t.Fatalf("authenticateWithAPIKey got err: %v", err)
	}
	want := "Successfully authenticated using the API key."
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("authenticateWithAPIKey got %q, want %q", got, want)
	}
}

func TestValidateServiceAccountKey(t *testing.T) {
	testutil.SystemTest(t)

	t.Run("valid key", func(t *testing.T) {
		keyPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if keyPath == "" {
			t.Skip("GOOGLE_APPLICATION_CREDENTIALS not set")
		}

		buf := &bytes.Buffer{}
		if err := validateServiceAccountKey(buf, keyPath); err != nil {
			t.Errorf("validateServiceAccountKey(valid) got err: %v", err)
		}

		want := "Successfully validated service account key"
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("validateServiceAccountKey(valid) got %q, want %q", got, want)
		}
	})

	t.Run("invalid key type", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "invalid-key-*.json")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		// A valid JSON but wrong type ("authorized_user" instead of "service_account").
		content := []byte(`{"type": "authorized_user"}`)
		if _, err := tmpFile.Write(content); err != nil {
			t.Fatal(err)
		}
		if err := tmpFile.Close(); err != nil {
			t.Fatalf("failed to close temp file: %v", err)
		}

		buf := &bytes.Buffer{}
		// The function should return an error because JWTConfigFromJSON
		// specifically expects the "service_account" type.
		if err := validateServiceAccountKey(buf, tmpFile.Name()); err == nil {
			t.Error("validateServiceAccountKey(invalid) expected error for 'authorized_user' type, got nil")
		}
	})
}
