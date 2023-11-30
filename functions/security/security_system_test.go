// Copyright 2020 Google LLC
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

package security

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const RuntimeVersion = "go119"
const FunctionsRegion = "us-central1"

func TestMain(m *testing.M) {
	// Only run end-to-end tests when configured to do so.
	if os.Getenv("GOLANG_SAMPLES_E2E_TEST") == "" {
		log.Println("Skipping end-to-end tests: GOLANG_SAMPLES_E2E_TEST not set")
		os.Exit(m.Run())
	}

	if os.Getenv("GOLANG_SAMPLES_PROJECT_ID") == "" {
		log.Println("Stopping test execution: GOLANG_SAMPLES_PROJECT_ID not set")
		os.Exit(0)
	}

	retn, err := setupAndRun(m)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(retn)
}

func setupAndRun(m *testing.M) (int, error) {
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	region := "us-central1"

	entryPoint := "MakeGetRequestCloudFunction"
	targetName := entryPoint + "-echo"

	// Setup function for tests.
	cmd := exec.Command("gcloud", "functions", "deploy", targetName,
		"--entry-point="+entryPoint,
		"--runtime="+RuntimeVersion,
		"--no-allow-unauthenticated",
		"--project="+os.Getenv("GOLANG_SAMPLES_PROJECT_ID"),
		"--trigger-http",
	)
	log.Printf("Running: %s %s", cmd.Path, strings.Join(cmd.Args, " "))
	if _, err := cmd.Output(); err != nil {
		log.Println(string(err.(*exec.ExitError).Stderr))
		return 1, fmt.Errorf("Setup: Deploy target function: %w", err)
	}
	defer teardown(targetName)

	// Setup function that tests directly use.
	targetURL := fmt.Sprintf("https://%s-%s.cloudfunctions.net/%s", region, projectID, targetName)
	cmd = exec.Command("gcloud", "functions", "deploy", entryPoint,
		"--entry-point="+entryPoint,
		"--runtime="+RuntimeVersion,
		"--no-allow-unauthenticated",
		"--update-env-vars", "TARGET_URL="+targetURL,
		"--project="+os.Getenv("GOLANG_SAMPLES_PROJECT_ID"),
		"--trigger-http",
	)
	log.Printf("Running: %s %s", cmd.Path, strings.Join(cmd.Args, " "))
	if _, err := cmd.Output(); err != nil {
		log.Println(string(err.(*exec.ExitError).Stderr))
		return 1, fmt.Errorf("Setup: Deploy relay function: %w", err)
	}
	defer teardown(entryPoint)

	baseURL := fmt.Sprintf("https://%s-%s.cloudfunctions.net/", region, projectID)
	os.Setenv("BASE_URL", baseURL)

	// Run the tests.
	return m.Run(), nil
}

func teardown(functionName string) {
	cmd := exec.Command("gcloud", "functions", "delete", functionName)
	log.Printf("Running: %s %s", cmd.Path, strings.Join(cmd.Args, " "))
	if _, err := cmd.Output(); err != nil {
		log.Println(string(err.(*exec.ExitError).Stderr))
		log.Printf("Teardown: Delete function %s: %v", functionName, err)
	}
}

func TestMakeGetRequest(t *testing.T) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		t.Skip("BASE_URL not set")
	}
	url := baseURL + "MakeGetRequestCloudFunction"

	testutil.Retry(t, 5, 30*time.Second, func(r *testutil.R) {
		var b bytes.Buffer
		if err := makeGetRequest(&b, url, url); err != nil {
			r.Errorf("makeGetRequest: %v", err)
		}
		got := b.String()
		if got != "Success!" {
			r.Errorf("got %s, want %s", got, "Success!")
		}
	})
}
