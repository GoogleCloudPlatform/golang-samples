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

package p

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		want     string
		wantCode int
	}{
		{
			name:     "valid",
			data:     `{"message": "Greetings, Ocean!"}`,
			want:     "Greetings, Ocean!",
			wantCode: http.StatusOK,
		},
		{
			name:     "empty",
			data:     "",
			want:     "Hello World!",
			wantCode: http.StatusOK,
		},
		{
			name:     "empty+braces",
			data:     "{}",
			want:     "Hello World!",
			wantCode: http.StatusOK,
		},
		{
			name:     "valid-no-message",
			data:     `{"data": "unused"}`,
			want:     "Hello World!",
			wantCode: http.StatusOK,
		},
		{
			name:     "invalid",
			data:     "not-valid-JSON",
			want:     http.StatusText(http.StatusBadRequest) + "\n",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("POST", "/", strings.NewReader(test.data))
		rr := httptest.NewRecorder()
		HelloWorld(rr, req)

		if got := rr.Result().StatusCode; got != test.wantCode {
			t.Errorf("HelloWorld(%s) Status: got '%d', want '%d'", test.name, got, test.wantCode)
		}

		if got := rr.Body.String(); got != test.want {
			t.Errorf("HelloWorld(%s) Body: got %q, want %q", test.name, got, test.want)
		}
	}
}
