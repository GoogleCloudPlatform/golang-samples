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
	"os"
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	if os.Getenv("GOLANG_SAMPLES_E2E_TEST") == "" {
		t.Skip()
	}
	tests := []struct {
		dbHost string
	}{
		{dbHost: os.Getenv("DB_HOST")},
	}

	for _, test := range tests {
		oldDBHost := os.Getenv("DB_HOST")
		os.Setenv("DB_HOST", test.dbHost)

		app := newApp()
		rr := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/", nil)
		app.indexHandler(rr, request)
		resp := rr.Result()
		body := rr.Body.String()

		if resp.StatusCode != 200 {
			t.Errorf("With dbHost='%s', indexHandler got status code %d, want 200", test.dbHost, resp.StatusCode)
		}

		want := "Tabs VS Spaces"
		if !strings.Contains(body, want) {
			t.Errorf("With dbHost='%s', expected to see '%s' in indexHandler response body", test.dbHost, want)
		}
		os.Setenv("DB_HOST", oldDBHost)
	}
}

func TestCastVote(t *testing.T) {
	if os.Getenv("GOLANG_SAMPLES_E2E_TEST") == "" {
		t.Skip()
	}
	tests := []struct {
		dbHost string
	}{
		{dbHost: os.Getenv("DB_HOST")},
	}

	for _, test := range tests {
		oldDBHost := os.Getenv("DB_HOST")
		os.Setenv("DB_HOST", test.dbHost)

		app := newApp()
		rr := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte("team=SPACES")))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		app.indexHandler(rr, request)
		resp := rr.Result()
		body := rr.Body.String()

		if resp.StatusCode != 200 {
			t.Errorf("With dbHost='%s', indexHandler got status code %d, want 200", test.dbHost, resp.StatusCode)
		}

		want := "Vote successfully cast for SPACES"
		if !strings.Contains(body, want) {
			t.Errorf("With dbHost='%s', expected to see '%s' in indexHandler response body", test.dbHost, want)
		}
		os.Setenv("DB_HOST", oldDBHost)
	}
}
