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
)

const RuntimeVersion = "go113"

func TestMain(m *testing.M) {
	// Only run end-to-end tests when configured to do so.
	if os.Getenv("GOLANG_SAMPLES_E2E_TEST") == "" {
		log.Println("Skipping end-to-end tests: GOLANG_SAMPLES_E2E_TEST not set")
		os.Exit(m.Run())
	}

	retn, err := setupAndRun(m)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(retn)
}

func setupAndRun(m *testing.M) (int, error) {
	entryPoint := "MakeGetRequestFixture"
	targetName := entryPoint + "-" + time.Now().Format("20060102-150405")

	// Setup function for tests.
	cmd := exec.Command("gcloud", "functions", "deploy", targetName,
		"--entry-point="+entryPoint,
		"--runtime="+RuntimeVersion,
		"--no-allow-unauthenticated",
		"--trigger-http",
	)
	log.Printf("Running: %s %s", cmd.Path, strings.Join(cmd.Args, " "))
	if _, err := cmd.Output(); err != nil {
		log.Println(string(err.(*exec.ExitError).Stderr))
		return 1, fmt.Errorf("Setup: Deploy target function: %w", err)
	}

	// Tear down the deployed function.
	defer func() {
		cmd = exec.Command("gcloud", "functions", "delete", targetName)
		log.Printf("Running: %s %s", cmd.Path, strings.Join(cmd.Args, " "))
		if _, err := cmd.Output(); err != nil {
			log.Println(string(err.(*exec.ExitError).Stderr))
			log.Printf("Teardown: Delete target function: %v", err)
		}
	}()

	// Setup Relay service where the authenticated request will be run.
	cmd = exec.Command("gcloud", "functions", "describe", targetName, "--format=value(httpsTrigger.url)")
	log.Printf("Running: %s %s", cmd.Path, strings.Join(cmd.Args, " "))
	out, err := cmd.Output()
	if err != nil {
		log.Println(string(err.(*exec.ExitError).Stderr))
		return 1, fmt.Errorf("Setup: Get target function URL: %w", err)
	}
	targetURL := strings.TrimSpace(string(out))
	relayName := entryPoint + "-relay-" + time.Now().Format("20060102-150405")

	// Setup function for tests.
	cmd = exec.Command("gcloud", "functions", "deploy", relayName,
		"--entry-point="+entryPoint,
		"--runtime="+RuntimeVersion,
		"--no-allow-unauthenticated",
		"--update-env-vars", "TARGET="+targetURL,
		"--trigger-http",
	)
	log.Printf("Running: %s %s", cmd.Path, strings.Join(cmd.Args, " "))
	if _, err := cmd.Output(); err != nil {
		log.Println(string(err.(*exec.ExitError).Stderr))
		return 1, fmt.Errorf("Setup: Deploy relay function: %w", err)
	}

	// Tear down the deployed function.
	defer func() {
		cmd = exec.Command("gcloud", "functions", "delete", relayName)
		log.Printf("Running: %s %s", cmd.Path, strings.Join(cmd.Args, " "))
		if _, err := cmd.Output(); err != nil {
			log.Println(string(err.(*exec.ExitError).Stderr))
			log.Printf("Teardown: Delete relay function: %v", err)
		}
	}()

	// Retrieve URL for tests.
	cmd = exec.Command("gcloud", "functions", "describe", relayName, "--format=value(httpsTrigger.url)")
	log.Printf("Running: %s %s", cmd.Path, strings.Join(cmd.Args, " "))
	out, err = cmd.Output()
	if err != nil {
		log.Println(string(err.(*exec.ExitError).Stderr))
		return 1, fmt.Errorf("Setup: Get relay function URL: %w", err)
	}
	if err := os.Setenv("BASE_URL", strings.TrimSpace(string(out))); err != nil {
		return 1, fmt.Errorf("Setup: os.Setenv: %w", err)
	}

	// Run the tests.
	return m.Run(), nil
}

func TestMakeGetRequest(t *testing.T) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		t.Skip("BASE_URL not set")
	}

	var b bytes.Buffer
	if err := makeGetRequest(&b, baseURL); err != nil {
		t.Fatalf("makeGetRequest: %v", err)
	}
	got := b.String()
	if got != "Success!" {
		t.Fatalf("got %s, want %s", got, "Success!")
	}
}
