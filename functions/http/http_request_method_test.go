// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package http

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHelloHTTPMethod(t *testing.T) {
	tests := []struct {
		method string
		want   string
	}{
		{method: "GET", want: "Hello World!"},
		{method: "PUT", want: "403 - Forbidden\n"},
		{method: "PATCH", want: "405 - Method Not Allowed\n"},
	}

	for _, test := range tests {
		payload := strings.NewReader("")

		req := httptest.NewRequest(test.method, "/", payload)
		rr := httptest.NewRecorder()

		HelloHTTPMethod(rr, req)
		out, err := ioutil.ReadAll(rr.Result().Body)
		if err != nil {
			t.Fatalf("%s: ReadAll: %v", test.method, err)
		}
		if got := string(out); got != test.want {
			t.Errorf("%s: got %q, want %q", test.method, got, test.want)
		}
	}
}
