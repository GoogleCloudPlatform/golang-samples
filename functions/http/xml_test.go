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

		if got := rr.Body.String(); got != test.want {
			t.Errorf("HelloHTTP(%q) = %q, want %q", test.body, got, test.want)
		}
	}
}
