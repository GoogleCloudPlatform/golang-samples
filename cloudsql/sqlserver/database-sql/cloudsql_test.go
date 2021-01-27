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

package main

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestIndex(t *testing.T) {
	testutil.EndToEndTest(t)
	tcpApp := createTCPApp()
	loadIndex(tcpApp, t)
}

func TestCastVote(t *testing.T) {
	testutil.EndToEndTest(t)
	tcpApp := createTCPApp()
	castVote(tcpApp, t)

}

func createTCPApp() *app {
	return newApp()
}

func loadIndex(app *app, t *testing.T) {
	rr := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	app.indexHandler(rr, request)
	resp := rr.Result()
	body := rr.Body.String()

	if resp.StatusCode != 200 {
		t.Errorf("indexHandler got status code %d, want 200", resp.StatusCode)
	}

	want := "Tabs VS Spaces"
	if !strings.Contains(body, want) {
		t.Errorf("Expected to see '%s' in indexHandler response body", want)
	}
}

func castVote(app *app, t *testing.T) {
	rr := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte("team=SPACES")))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	app.indexHandler(rr, request)
	resp := rr.Result()
	body := rr.Body.String()

	if resp.StatusCode != 200 {
		t.Errorf("indexHandler got status code %d, want 200", resp.StatusCode)
	}

	want := "Vote successfully cast for SPACES"
	if !strings.Contains(body, want) {
		t.Errorf("Expected to see '%s' in indexHandler response body", want)
	}
}
