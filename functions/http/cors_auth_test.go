// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCORSEnabledFunctionAuth(t *testing.T) {
	req := httptest.NewRequest("OPTIONS", "/", strings.NewReader(""))

	rr := httptest.NewRecorder()
	CORSEnabledFunctionAuth(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("CORSEnabledFunction got status %v, want %v", rr.Code, http.StatusNoContent)
	}
	headers := map[string]string{
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Headers":     "Authorization",
		"Access-Control-Allow-Methods":     "POST",
		"Access-Control-Allow-Origin":      "https://example.com",
		"Access-Control-Max-Age":           "3600",
	}
	for k, v := range headers {
		if got := rr.Header().Get(k); got != v {
			t.Errorf("CORSEnabledFunctionAuth header[%v] = %v, want %v", k, got, v)
		}
	}
}

func TestCORSEnabledFunctionAuthPOST(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader(""))

	rr := httptest.NewRecorder()
	CORSEnabledFunctionAuth(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("CORSEnabledFunction got status %v, want %v", rr.Code, http.StatusOK)
	}

	headers := map[string]string{
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Origin":      "https://example.com",
	}
	for k, v := range headers {
		if got := rr.Header().Get(k); got != v {
			t.Errorf("CORSEnabledFunctionAuth header[%v] = %v, want %v", k, got, v)
		}
	}
}

func TestCORSEnabledFunctionAuthSystem(t *testing.T) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		t.Skip("BASE_URL not set")
	}
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	urlString := baseURL + "/CORSEnabledFunctionAuth"
	testURL, err := url.Parse(urlString)
	if err != nil {
		t.Fatalf("url.Parse(%q): %v", urlString, err)
	}

	req := &http.Request{
		Method: http.MethodOptions,
		Body:   ioutil.NopCloser(strings.NewReader("")),
		URL:    testURL,
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("HelloHTTP http.Get: %v", err)
	}

	headers := map[string]string{
		"Access-Control-Allow-Origin":  "https://example.com",
		"Access-Control-Allow-Methods": "POST",
		"Access-Control-Allow-Headers": "Authorization",
		"Access-Control-Max-Age":       "3600",
	}
	for k, v := range headers {
		if got := resp.Header.Get(k); got != v {
			t.Errorf("CORSEnabledFunctionAuth header[%v] = %v, want %v", k, got, v)
		}
	}
}
