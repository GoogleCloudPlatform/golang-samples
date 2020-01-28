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

func TestIndexHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", strings.NewReader(""))
	rr := httptest.NewRecorder()
	indexHandler(rr, req)

	if got := rr.Result().StatusCode; got != http.StatusOK {
		t.Errorf("response status: got %q, want %q", got, http.StatusOK)
	}

	want := `<title>Markdown Editor</title>`
	re := regexp.MustCompile(`<title>.*</title>`)
	got := re.FindString(rr.Body.String())

	if got != want {
		t.Errorf("body: got %q, want %q", got, want)
	}
}
