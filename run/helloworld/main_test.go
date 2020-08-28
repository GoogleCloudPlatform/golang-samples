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
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		label string
		want  string
		name  string
	}{
		{
			label: "default",
			want:  "Hello World!\n",
			name:  "",
		},
		{
			label: "override",
			want:  "Hello Override!\n",
			name:  "Override",
		},
	}

	originalName := os.Getenv("NAME")
	defer os.Setenv("NAME", originalName)

	for _, test := range tests {
		os.Setenv("NAME", test.name)

		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()
		handler(rr, req)

		if got := rr.Body.String(); got != test.want {
			t.Errorf("%s: got %q, want %q", test.label, got, test.want)
		}
	}
}
