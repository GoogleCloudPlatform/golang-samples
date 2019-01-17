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

func TestParseXML(t *testing.T) {
	tests := []struct {
		body string
		want string
	}{
		{
			body: `<Person><Name>Gopher</Name></Person>`,
			want: "Hello, Gopher!",
		},
		{
			body: `<Person></Person>`,
			want: "Hello, World!",
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/", strings.NewReader(test.body))

		rr := httptest.NewRecorder()
		ParseXML(rr, req)

		out, err := ioutil.ReadAll(rr.Result().Body)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}
		if got := string(out); got != test.want {
			t.Errorf("HelloHTTP(%q) = %q, want %q", test.body, got, test.want)
		}
	}
}
