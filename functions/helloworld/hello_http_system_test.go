// Copyright 2019 Google LLC
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

//go:build ignore
// +build ignore

// Disabled until system tests are working on Kokoro.

// [START functions_http_system_test]

package helloworld

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
	entryPoint := "HelloHTTP"
	name := entryPoint + "-" + time.Now().Format("20060102-150405")

	// Setup function for tests.
	cmd := exec.Command("gcloud", "functions", "deploy", name,
		"--entry-point="+entryPoint,
		"--runtime="+RuntimeVersion,
		"--allow-unauthenticated",
		"--trigger-http",
	)
	log.Printf("Running: %s %s", cmd.Path, strings.Join(cmd.Args, " "))
	if _, err := cmd.Output(); err != nil {
		log.Println(string(err.(*exec.ExitError).Stderr))
		return 1, fmt.Errorf("Setup: Deploy function: %w", err)
	}

	// Tear down the deployed function.
	defer func() {
		cmd = exec.Command("gcloud", "functions", "delete", name)
		log.Printf("Running: %s %s", cmd.Path, strings.Join(cmd.Args, " "))
		if _, err := cmd.Output(); err != nil {
			log.Println(string(err.(*exec.ExitError).Stderr))
			log.Printf("Teardown: Delete function: %v", err)
		}
	}()

	// Retrieve URL for tests.
	cmd = exec.Command("gcloud", "functions", "describe", name, "--format=value(httpsTrigger.url)")
	log.Printf("Running: %s %s", cmd.Path, strings.Join(cmd.Args, " "))
	out, err := cmd.Output()
	if err != nil {
		log.Println(string(err.(*exec.ExitError).Stderr))
		return 1, fmt.Errorf("Setup: Get function URL: %w", err)
	}
	if err := os.Setenv("BASE_URL", strings.TrimSpace(string(out))); err != nil {
		return 1, fmt.Errorf("Setup: os.Setenv: %w", err)
	}

	// Run the tests.
	return m.Run(), nil
}

func TestHelloHTTPSystem(t *testing.T) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	urlString := os.Getenv("BASE_URL")
	testURL, err := url.Parse(urlString)
	if err != nil {
		t.Fatalf("url.Parse(%q): %v", urlString, err)
	}

	tests := []struct {
		body string
		want string
	}{
		{body: `{"name": ""}`, want: "Hello, World!"},
		{body: `{"name": "Gopher"}`, want: "Hello, Gopher!"},
	}

	for _, test := range tests {
		req := &http.Request{
			Method: http.MethodPost,
			Body:   io.NopCloser(strings.NewReader(test.body)),
			URL:    testURL,
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("HelloHTTP http.Get: %v", err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("HelloHTTP io.ReadAll: %v", err)
		}
		if got := string(body); got != test.want {
			t.Errorf("HelloHTTP(%q) = %q, want %q", test.body, got, test.want)
		}
	}
}

// [END functions_http_system_test]
