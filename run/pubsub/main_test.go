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
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHelloPubSubErrors(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "no_payload",
			message: "",
		},
		{
			name:    "not_base64",
			message: `{"message":{"data":"Gopher","id":"test-123"}}`,
		},
	}
	for _, test := range tests {
		payload := strings.NewReader(test.message)
		req := httptest.NewRequest("GET", "/", payload)
		rr := httptest.NewRecorder()

		HelloPubSub(rr, req)

		if code := rr.Result().StatusCode; code != http.StatusBadRequest {
			t.Errorf("HelloPubSub(%q): got (%q), want (%q)", test.name, code, http.StatusBadRequest)
		}
	}
}

func TestHelloPubSub(t *testing.T) {
	tests := []struct {
		data string
		want string
	}{
		{want: "Hello World!\n"},
		{data: "Go", want: "Hello Go!\n"},
	}
	for _, test := range tests {
		r, w, _ := os.Pipe()
		log.SetOutput(w)
		originalFlags := log.Flags()
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

		payload := strings.NewReader("{}")
		if test.data != "" {
			encoded := base64.StdEncoding.EncodeToString([]byte(test.data))
			jsonStr := fmt.Sprintf(`{"message":{"data":"%s","id":"test-123"}}`, encoded)
			payload = strings.NewReader(jsonStr)
		}
		req := httptest.NewRequest("GET", "/", payload)
		rr := httptest.NewRecorder()

		HelloPubSub(rr, req)

		w.Close()
		log.SetOutput(os.Stderr)
		log.SetFlags(originalFlags)

		if code := rr.Result().StatusCode; code == http.StatusBadRequest {
			t.Errorf("HelloPubSub(%q) invalid input, status code (%q)", test.data, code)
		}

		out, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}
		if got := string(out); got != test.want {
			t.Errorf("HelloPubSub(%q): got %q, want %q", test.data, got, test.want)
		}
	}
}
