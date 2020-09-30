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

package main_test

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	cr "github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// TestApp shows the simplest approach to create e2e tests with cloudrun-ci.
// It uses Cloud Run (fully managed) in the us-central1 region, and will manage
// container images as needed for the test.
func TestApp(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	helloworld := cr.NewService("hello", tc.ProjectID)
	if err := helloworld.Deploy(); err != nil {
		t.Fatalf("could not deploy %s: %v", helloworld.Name, err)
	}
	defer helloworld.Clean()

	bodyShouldContain(t, helloworld, "/", "Hello World!")
}

func bodyShouldContain(t *testing.T, s *cr.Service, path, shouldContain string) {
	if !s.Deployed() {
		t.Fatalf("[%s] not deployed", s.Name)
	}

	url, err := s.URL("")
	if err != nil {
		t.Fatalf("[%s] could not retrieve URL", s.Name)
	}
	log.Printf("[%s] Deployed to %s", s.Name, url)

	resp, err := s.Request("GET", path)
	if err != nil {
		t.Errorf("Request: %v", err)
		return
	}
	defer resp.Body.Close()

	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll: %q", err)
	}

	if got := string(got); !strings.Contains(got, shouldContain) {
		t.Errorf("got\n----\n%s\n----\nWant to contain:\n----\n%s\n", got, shouldContain)
	}
}
