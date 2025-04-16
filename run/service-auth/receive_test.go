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
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

func TestAuth(t *testing.T) {
	tc := testutil.SystemTest(t)

	serviceName := testGenerateServiceID()

	if err := testDeployReceiveService(t, serviceName, tc.ProjectID); err != nil {
		t.Fatalf("testDeployReceiveService error: %v\n", err)
	}
	defer testDeleteReceiveService(t, serviceName, tc.ProjectID)

	url, err := testGetReceiveServiceURL(t, serviceName, tc.ProjectID)
	if err != nil {
		t.Fatalf("testGetReceiveServiceURL error: %v\n", err)
	}

	token, err := testGetGCPAuthToken(t)
	if err != nil {
		t.Fatalf("testGetGCPAuthToken error: %v\n", err)
	}

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		request, err := http.NewRequest(http.MethodGet, url, nil)
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

func testGetGCPAuthToken(t *testing.T) (string, error) {
	t.Helper()

	cmd := exec.Command(
		"gcloud",
		"auth",
		"print-identity-token",
	)

	bytesToken, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("cmd.Run error: %w", err)
	}

	token := strings.ReplaceAll(string(bytesToken), "\n", "")

	return token, nil
}

func testGetReceiveServiceURL(t *testing.T, serviceName, projectID string) (string, error) {
	t.Helper()

	cmd := exec.Command(
		"gcloud",
		"run",
		"services",
		"describe",
		serviceName,
		"--project", projectID,
		"--region=us-central1",
		"--format",
		"value(status.url)",
	)

	bytesURL, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("cmd.Run error: %w", err)
	}

	url := strings.ReplaceAll(string(bytesURL), "\n", "")

	return url, nil
}

func testDeployReceiveService(t *testing.T, serviceName, projectID string) error {
	t.Helper()

	cmd := exec.Command(
		"gcloud",
		"run",
		"deploy",
		serviceName,
		"--project="+projectID,
		"--source", ".",
		"--region=us-central1",
		"--allow-unauthenticated",
		"--quiet",
	)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cmd.Run error: %w", err)
	}

	return nil
}

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

func testGenerateServiceID() string {
	return fmt.Sprintf("receive-go-%s", uuid.New().String())
}
