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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestStartingServer(t *testing.T) {
	// tc := testutil.SystemTest(t)
	// test build
	m := testutil.BuildMain(t)
	if !m.Built() {
		t.Fatalf("failed to build app")
	}

	// test main endpoint
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle)
	handler.ServeHTTP(rr, req)

	// test response status
	gotStatus := rr.Code
	wantStatus := http.StatusOK
	fmt.Print("status was: " + strconv.Itoa(gotStatus) + "\n")

	if gotStatus != wantStatus {
		t.Fatalf("Returned wrong status code: got %v want %v",
			gotStatus, wantStatus)
	}

	// test response body
	// acceptable responses are "intentional error!" and "succeeded after ..."
	wantSuccess := "Succeeded after "
	wantIntentionalError := "intentional error!"

	if rr.Body.String() != wantIntentionalError && !strings.Contains(rr.Body.String(), wantSuccess) {
		t.Fatalf("Response does not match expected: got %v",
			rr.Body.String())
	}

}
