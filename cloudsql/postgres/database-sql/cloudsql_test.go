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

type testInfo struct {
	dbName                 string
	dbPass                 string
	dbUser                 string
	dbPort                 string
	instanceConnectionName string
}

func TestIndex(t *testing.T) {
	if os.Getenv("GOLANG_SAMPLES_E2E_TEST") == "" {
		t.Skip()
	}

	info := testInfo{
		dbName:                 os.Getenv("POSTGRES_DATABASE"),
		dbPass:                 os.Getenv("POSTGRES_PASSWORD"),
		dbPort:                 os.Getenv("POSTGRES_PORT"),
		dbUser:                 os.Getenv("POSTGRES_USER"),
		instanceConnectionName: os.Getenv("POSTGRES_INSTANCE"),
	}

	tests := []struct {
		dbHost string
	}{
		{dbHost: ""},
		{dbHost: os.Getenv("POSTGRES_HOST")},
	}

	// Capture original values
	oldDBHost := os.Getenv("DB_HOST")
	oldDBName := os.Getenv("DB_NAME")
	oldDBPass := os.Getenv("DB_PASS")
	oldDBPort := os.Getenv("DB_PORT")
	oldDBUser := os.Getenv("DB_USER")
	oldInstance := os.Getenv("INSTANCE_CONNECTION_NAME")

	for _, test := range tests {

		// Set overwrites
		os.Setenv("DB_HOST", test.dbHost)
		os.Setenv("DB_NAME", info.dbName)
		os.Setenv("DB_PASS", info.dbPass)
		os.Setenv("DB_PORT", info.dbPort)
		os.Setenv("DB_USER", info.dbUser)
		os.Setenv("INSTANCE_CONNECTION_NAME", info.instanceConnectionName)

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

	}
	// Restore original values
	os.Setenv("DB_HOST", oldDBHost)
	os.Setenv("DB_NAME", oldDBName)
	os.Setenv("DB_PASS", oldDBPass)
	os.Setenv("DB_PORT", oldDBPort)
	os.Setenv("DB_USER", oldDBUser)
	os.Setenv("INSTANCE_CONNECTION_NAME", oldInstance)
}

func TestCastVote(t *testing.T) {
	if os.Getenv("GOLANG_SAMPLES_E2E_TEST") == "" {
		t.Skip()
	}

	info := testInfo{
		dbName:                 os.Getenv("POSTGRES_DATABASE"),
		dbPass:                 os.Getenv("POSTGRES_PASSWORD"),
		dbPort:                 os.Getenv("POSTGRES_PORT"),
		dbUser:                 os.Getenv("POSTGRES_USER"),
		instanceConnectionName: os.Getenv("POSTGRES_INSTANCE"),
	}

	tests := []struct {
		dbHost string
	}{
		{dbHost: ""},
		{dbHost: os.Getenv("POSTGRES_HOST")},
	}

	// Capture original values
	oldDBHost := os.Getenv("DB_HOST")
	oldDBName := os.Getenv("DB_NAME")
	oldDBPass := os.Getenv("DB_PASS")
	oldDBPort := os.Getenv("DB_PORT")
	oldDBUser := os.Getenv("DB_USER")
	oldInstance := os.Getenv("INSTANCE_CONNECTION_NAME")

	for _, test := range tests {

		// Set overwrites
		os.Setenv("DB_HOST", test.dbHost)
		os.Setenv("DB_NAME", info.dbName)
		os.Setenv("DB_PASS", info.dbPass)
		os.Setenv("DB_PORT", info.dbPort)
		os.Setenv("DB_USER", info.dbUser)
		os.Setenv("INSTANCE_CONNECTION_NAME", info.instanceConnectionName)

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

	// Restore original values
	os.Setenv("DB_HOST", oldDBHost)
	os.Setenv("DB_NAME", oldDBName)
	os.Setenv("DB_PASS", oldDBPass)
	os.Setenv("DB_PORT", oldDBPort)
	os.Setenv("DB_USER", oldDBUser)
	os.Setenv("INSTANCE_CONNECTION_NAME", oldInstance)
}
