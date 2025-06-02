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
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

const REGION string = "us-central1"

func TestAuth(t *testing.T) {
	tc := testutil.SystemTest(t)

	// Generate Service ID
	serviceName := testGenerateServiceID()

	// Once the test is done, the created resources are released
	defer testDeleteReceiveService(t, serviceName, tc.ProjectID)

	// Get ProjectNumber
	projectNumber, err := testGetProjectNumber(t, tc.ProjectID)
	if err != nil {
		t.Fatalf("testGetProjectNumber error: %v\n", err)
	}

	// Build serviceURL with expected format.
	serviceURL := fmt.Sprintf("https://%s-%s.%s.run.app", serviceName, projectNumber, REGION)

	// Deploy service
	if err := testDeployReceiveService(t, serviceName, tc.ProjectID, serviceURL); err != nil {
		t.Fatalf("testDeployReceiveService error: %v\n", err)
	}

	// Generate authentication token.
	token, err := testGetGCPAuthToken(t, serviceURL)
	if err != nil {
		t.Fatalf("testGetGCPAuthToken error: %v\n", err)
	}

	// Test deployed service with retry configuration.
	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		request, err := http.NewRequest(http.MethodGet, serviceURL, nil)
		if err != nil {
			r.Errorf("Attempt %v http.NewRequest error: %v", r.Attempt, err)
		}

		request.Header.Set("Authorization", "Bearer "+token)

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			r.Errorf("Attempt %v http.DefaultClient.Do error: %v", r.Attempt, err)
		}
		defer response.Body.Close()

		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			r.Errorf("Attempt %v io.ReadAll error: %v", r.Attempt, err)
		}

		if got, want := response.StatusCode, http.StatusOK; got != want {
			r.Errorf("Attempt %v Receive Service: unexpected status got %v want %v.\n", r.Attempt, got, want)
		}

		if got, dontWant := string(responseBody), "anonymous"; strings.Contains(got, dontWant) {
			r.Errorf("Attempt %v Receive Service: got: %s dont want %q\n", r.Attempt, got, dontWant)
		}
	})
}

// testGetGCPAuthToken returns an access token with specified audience.
func testGetGCPAuthToken(t *testing.T, endpointURL string) (string, error) {
	t.Helper()
	ctx := context.Background()

	// idtoken.NewTokenSource creates a TokenSource that can provide ID tokens
	// for the given audience. It handles the underlying fetching and refreshing.
	tokenSource, err := idtoken.NewTokenSource(ctx, endpointURL)
	if err != nil {
		return "", fmt.Errorf("idtoken.NewTokenSource: %w", err)
	}

	// Call Token() on the tokenSource to get the actual ID token.
	idToken, err := tokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("tokenSource.Token: %w", err)
	}

	// The ID token string is in the AccessToken field of the oauth2.Token struct.
	return idToken.AccessToken, nil
}

// testGetProjectNumber returns formatted Project Number with given projectID.
func testGetProjectNumber(t *testing.T, projectID string) (string, error) {
	t.Helper()

	cmd := exec.Command(
		"gcloud",
		"projects",
		"describe",
		projectID,
		"--format=value(projectNumber)",
	)

	bytesURL, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("cmd.Run error: %w", err)
	}

	projectNumber := strings.ReplaceAll(string(bytesURL), "\n", "")

	return projectNumber, nil
}

// testDeployReceiveService deploys current source code to CloudRun.
func testDeployReceiveService(t *testing.T, serviceName, projectID, serviceURL string) error {
	t.Helper()

	cmd := exec.Command(
		"gcloud",
		"run",
		"deploy",
		serviceName,
		"--project",
		projectID,
		"--source",
		".",
		"--region="+REGION,
		"--allow-unauthenticated",
		"--set-env-vars=SERVICE_URL="+serviceURL,
		"--quiet",
	)

	cmd.Stdout = os.Stdout // or any other io.Writer
	cmd.Stderr = os.Stderr // or any other io.Writer

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cmd.Run error: %w", err)
	}

	return nil
}

// testDeleteReceiveService deletes a deployed project in CloudRun with the given serviceName and projectID.
func testDeleteReceiveService(t *testing.T, serviceName, projectID string) error {
	t.Helper()

	cmd := exec.Command(
		"gcloud",
		"run",
		"services",
		"delete",
		serviceName,
		"--project="+projectID,
		"--async",
		"--region=us-central1",
		"--quiet",
	)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cmd.Run error: %w", err)
	}

	return nil
}

// testGenerateServiceID produces an unique ID.
func testGenerateServiceID() string {
	return fmt.Sprintf("receive-go-%s", uuid.New().String())
}
