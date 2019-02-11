// Copyright 2019 Google LLC
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
