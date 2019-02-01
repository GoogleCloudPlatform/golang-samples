// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_http_unit_test]

package helloworld

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHelloHTTP(t *testing.T) {
	tests := []struct {
		body string
		want string
	}{
		{body: `{"name": ""}`, want: "Hello, World!"},
		{body: `{"name": "Gopher"}`, want: "Hello, Gopher!"},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/", strings.NewReader(test.body))
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		HelloHTTP(rr, req)

		out, err := ioutil.ReadAll(rr.Result().Body)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}
		if got := string(out); got != test.want {
			t.Errorf("HelloHTTP(%q) = %q, want %q", test.body, got, test.want)
		}
	}
}

// [END functions_http_unit_test]
