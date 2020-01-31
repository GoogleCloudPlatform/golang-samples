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
	"os"
	"strings"
	"testing"
)

func TestBrokenErrors(t *testing.T) {
	payload := strings.NewReader("")
	req := httptest.NewRequest("GET", "/", payload)
	rr := httptest.NewRecorder()

	os.Setenv("NAME", "")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("helloHandler: got (no panic), want (panic)")
		}
	}()
	helloHandler(rr, req)
}

func TestBrokenHandler(t *testing.T) {
	tests := []struct {
		label string
		name  string
		want  string
	}{
		{
			label: "<SET>",
			name:  "Testers",
			want:  "Hello Testers!\n",
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/", strings.NewReader(""))
		rr := httptest.NewRecorder()

		os.Setenv("NAME", test.name)
		improvedHandler(rr, req)

		if code := rr.Result().StatusCode; code != http.StatusOK {
			t.Errorf("brokenHandler(%s): got (%q), want (%q)", test.label, code, http.StatusOK)
		}

		if got := rr.Body.String(); test.want != got {
			t.Errorf("brokenHandler(%s): got (%q), want (%q)", test.label, got, test.want)
		}
	}
}

func TestImprovedHandler(t *testing.T) {
	tests := []struct {
		label string
		name  string
		want  string
	}{
		{
			label: "<EMPTY>",
			name:  "",
			want:  "Hello World!\n",
		},
		{
			label: "<SET>",
			name:  "Testers",
			want:  "Hello Testers!\n",
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/", strings.NewReader(""))
		rr := httptest.NewRecorder()

		os.Setenv("NAME", test.name)
		improvedHandler(rr, req)

		if code := rr.Result().StatusCode; code != http.StatusOK {
			t.Errorf("brokenHandler(%s): got (%q), want (%q)", test.label, code, http.StatusOK)
		}

		if got := rr.Body.String(); test.want != got {
			t.Errorf("brokenHandler(%s): got (%q), want (%q)", test.label, got, test.want)
		}
	}
}
