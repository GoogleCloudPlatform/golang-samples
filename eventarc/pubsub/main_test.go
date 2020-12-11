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
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHelloPubSubCloudEvent(t *testing.T) {
	tests := []struct {
		data string
		want string
		id   string
	}{
		{want: "Hello, World! ID: \n", id: ""},
		{want: "Hello, World! ID: 12345\n", id: "12345"},
		{data: "Go", want: "Hello, Go! ID: \n"},
		{data: "Go", want: "Hello, Go! ID: 1234\n", id: "1234"},
	}
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	for _, test := range tests {
		r, w, _ := os.Pipe()
		log.SetOutput(w)
		defer log.SetOutput(os.Stderr)

		payload := strings.NewReader("{}")
		if test.data != "" {
			encoded := base64.StdEncoding.EncodeToString([]byte(test.data))
			jsonStr := fmt.Sprintf(`{"message":{"data":"%s","id":"%s"}}`, encoded, test.id)
			payload = strings.NewReader(jsonStr)
		}

		req := httptest.NewRequest("POST", "/", payload)
		req.Header.Set("Ce-Id", test.id)
		rr := httptest.NewRecorder()
		HelloEventsPubSub(rr, req)

		w.Close()

		if code := rr.Result().StatusCode; code == http.StatusBadRequest {
			t.Errorf("HelloEventsPubSub(%q) invalid input, status code (%q)", test.data, code)
		}

		out, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}
		if got := string(out); got != test.want {
			t.Errorf("HelloEventsPubSub(%q): got %q, want %q", test.data, got, test.want)
		}
	}
}
