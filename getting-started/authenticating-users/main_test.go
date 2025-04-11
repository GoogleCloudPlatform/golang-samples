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

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndex(t *testing.T) {
	tests := []struct {
		path           string
		wantStatusCode int
		wantBody       string
	}{
		{
			path:           "/",
			wantStatusCode: http.StatusOK,
			wantBody:       "No Cloud IAP header found.\n",
		},
		{
			path:           "/hello",
			wantStatusCode: http.StatusNotFound,
			wantBody:       "404 page not found\n",
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", test.path, nil)
		rr := httptest.NewRecorder()

		a := &app{} // Do not use newApp since it uses the metadata server.
		a.index(rr, req)

		if got := rr.Result().StatusCode; got != test.wantStatusCode {
			t.Errorf("index(%s) got status code %d, want %d", test.path, got, test.wantStatusCode)
		}

		if got := rr.Body.String(); got != test.wantBody {
			t.Errorf("index(%s) got %q, want %q", test.path, got, test.wantBody)
		}
	}
}
