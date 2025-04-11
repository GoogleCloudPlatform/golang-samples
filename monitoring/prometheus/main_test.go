// Copyright 2021 Google LLC
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
	"strconv"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestStartingServer(t *testing.T) {

	// Tests build.
	m := testutil.BuildMain(t)
	if !m.Built() {
		t.Fatalf("failed to build app")
	}

	// Tests main endpoint.
	req := httptest.NewRequest("GET", "/", strings.NewReader(""))
	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handle(rr, req)

	// Tests response status.
	gotStatus := rr.Code
	gotBody := rr.Body.String()
	wantGoodStatus := http.StatusOK
	wantBadStatus := http.StatusInternalServerError

	if gotStatus != wantGoodStatus && gotStatus != wantBadStatus {
		t.Fatalf("Returned wrong status code: got %v want %v", gotStatus, strconv.Itoa(wantGoodStatus)+" or "+strconv.Itoa(wantBadStatus))
	}

	// Test response body.
	// Acceptable responses are "intentional error!" and "succeeded after ...".
	wantSuccess := "Succeeded after "
	wantIntentionalError := "intentional error!"

	if gotBody != wantIntentionalError && !strings.Contains(gotBody, wantSuccess) {
		t.Fatalf("Response does not match expected: got %v", gotBody)
	}

}

func TestMetricsEndpoint(t *testing.T) {

	// Tests build.
	m := testutil.BuildMain(t)
	if !m.Built() {
		t.Fatalf("failed to build app")
	}

	// Tests metrics endpoint.
	req := httptest.NewRequest("GET", "/metrics", nil)
	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := promhttp.Handler()
	handler.ServeHTTP(rr, req)

	// Tests response status.
	gotStatus := rr.Code
	gotBody := rr.Body.String()
	wantGoodStatus := http.StatusOK

	if gotStatus != wantGoodStatus {
		t.Fatalf("Returned wrong status code: got %v want %v", gotStatus, wantGoodStatus)
	}

	// Test response body.
	// Prometheus /metrics endpoint will contain "# HELP".
	wantSuccess := "# HELP"

	if !strings.Contains(gotBody, wantSuccess) {
		t.Fatalf("Response does not match expected: got %v", rr.Body.String())
	}

}
