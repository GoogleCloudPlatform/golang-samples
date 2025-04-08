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
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status: got (%v) want (%v)", status, http.StatusOK)
	}

	want := "Hello, World!"
	if got := rr.Body.String(); !strings.Contains(got, want) {
		t.Errorf("unexpected body: got (%v) want (%v)", got, want)
	}
}

func TestSetup(t *testing.T) {
	setup(context.Background())
	if startupTime.IsZero() {
		t.Error("warmupApp: got (uninitialized startupTime) want (startupTime)")
	}
	if client == nil {
		t.Error("warmupApp: got (uninitialized client) want (client)")
	}
}
