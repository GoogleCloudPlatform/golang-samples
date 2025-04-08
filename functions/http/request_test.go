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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMakeRequest(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
	defer s.Close()
	// Override the url global with the test server URL.
	urlString = s.URL

	req := httptest.NewRequest("GET", "/", strings.NewReader(""))
	rr := httptest.NewRecorder()
	MakeRequest(rr, req)

	if got, want := rr.Body.String(), "ok"; got != want {
		t.Errorf("MakeRequest got %q, want %q", got, want)
	}
}
