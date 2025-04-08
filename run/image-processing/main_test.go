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
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHelloPubSubErrors(t *testing.T) {
	tests := []struct {
		name    string
		message string
		data    string
		encode  bool
	}{
		{
			name:    "no_payload",
			message: "",
		},
		{
			name:    "not_base64_invalid_props",
			message: `{"message":{"data":"Gopher","id":"test-123"}}`,
		},
		{
			name:   "no_name",
			data:   `{"bucket":"my-bucket"}`,
			encode: true,
		},
		{
			name:   "no_bucket",
			data:   `{"name":"my-object"}`,
			encode: true,
		},
		{
			name:    "not_base64_valid_props",
			message: `{"name":"my-object","bucket":"my-bucket" }`,
			encode:  false,
		},
	}
	for _, test := range tests {
		if test.message == "" && test.data != "" {
			data := test.data
			if test.encode {
				data = base64.StdEncoding.EncodeToString([]byte(data))
			}

			test.message = fmt.Sprintf(`{"message": {"data": "%s"}}`, data)
		}

		payload := strings.NewReader(test.message)
		req := httptest.NewRequest("POST", "/", payload)
		rr := httptest.NewRecorder()

		HelloPubSub(rr, req)

		if code := rr.Result().StatusCode; code != http.StatusBadRequest {
			t.Errorf("HelloPubSub(%q): got (%q), want (%q)", test.name, code, http.StatusBadRequest)
		}
	}
}
