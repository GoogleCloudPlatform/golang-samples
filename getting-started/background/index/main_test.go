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
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestIndex(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("Skipping Firestore test. Set GOLANG_SAMPLES_FIRESTORE_PROJECT.")
	}

	a, err := newApp(projectID, "")
	if err != nil {
		t.Fatalf("newApp: %v", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.index(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("wrong status code, got %v, want %v", resp.StatusCode, http.StatusOK)
	}
}
