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

package main

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("GOLANG_SAMPLES_FIRESTORE_PROJECT not set")
	}

	a, err := newApp(projectID)
	if err != nil {
		t.Fatalf("newApp: %v", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	a.index(rr, r)

	if got, want := rr.Body.String(), "1 view"; !strings.Contains(got, want) {
		t.Errorf("index first visit got:\n----\n%v\n----\nWant to contain %q", got, want)
	}

	// Subsequent requests include the cookie from first visit.
	r = httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Cookie", rr.Header().Get("Set-Cookie"))

	rr = httptest.NewRecorder()

	a.index(rr, r)

	if got, want := rr.Body.String(), "2 views"; !strings.Contains(got, want) {
		t.Errorf("index second visit got:\n----\n%v\n----\nWant to contain %q", got, want)
	}
}

func TestIndexCorrupted(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("GOLANG_SAMPLES_FIRESTORE_PROJECT not set")
	}

	a, err := newApp(projectID)
	if err != nil {
		t.Fatalf("newApp: %v", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Cookie", "this is not a valid session ID")

	rr := httptest.NewRecorder()

	a.index(rr, r)

	if got, want := rr.Body.String(), "1 view"; !strings.Contains(got, want) {
		t.Errorf("index first visit got:\n----\n%v\n----\nWant to contain %q", got, want)
	}
}
