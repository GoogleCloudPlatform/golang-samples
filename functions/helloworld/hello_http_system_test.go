// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// +build ignore
// Disabled until system tests are working on Kokoro.

// [START functions_http_system_test]

package helloworld

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHelloHTTPSystem(t *testing.T) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	urlString := os.Getenv("BASE_URL") + "/HelloHTTP"
	testURL, err := url.Parse(urlString)
	if err != nil {
		t.Fatalf("url.Parse(%q): %v", urlString, err)
	}

	tests := []struct {
		body string
		want string
	}{
		{body: `{"name": ""}`, want: "Hello, World!"},
		{body: `{"name": "Gopher"}`, want: "Hello, Gopher!"},
	}

	for _, test := range tests {
		req := &http.Request{
			Method: http.MethodPost,
			Body:   ioutil.NopCloser(strings.NewReader(test.body)),
			URL:    testURL,
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("HelloHTTP http.Get: %v", err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("HelloHTTP ioutil.ReadAll: %v", err)
		}
		if got := string(body); got != test.want {
			t.Errorf("HelloHTTP(%q) = %q, want %q", test.body, got, test.want)
		}
	}
}

// [END functions_http_system_test]
