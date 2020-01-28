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
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestRenderHandlerErrors(t *testing.T) {
	tests := []struct {
		label string
		req *http.Request
		wantBody string
		wantStatus int
	}{
		{
			label: "Invalid Method",
			req: httptest.NewRequest("GET", "/render", strings.NewReader("")),
			wantBody: http.StatusText(http.StatusMethodNotAllowed),
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			label: "Invalid JSON",
			req: httptest.NewRequest("POST", "/render", strings.NewReader("**markdown**")),
			wantBody: http.StatusText(http.StatusBadRequest),
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		rr := httptest.NewRecorder()
		indexHandler(rr, test.req)

		if got := rr.Result().StatusCode; got != test.wantStatus {
			t.Errorf("%s: response status: got %q, want %q", test.label, got, test.wantStatus)
		}

		if got := rr.Body.String(); got != test.wantBody {
			t.Errorf("%s: body: got %q, want %q", test.label, got, test.wantBody)
		}
	}
}
