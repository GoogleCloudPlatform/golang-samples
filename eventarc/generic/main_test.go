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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGenericCloudEvent(t *testing.T) {
	tests := []struct {
		data string
		want string
	}{
		{want: "Event received!"},
		{want: `"Ce-Id": "1234"`},
		{want: `{"message": "some string"}`},
	}
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	for _, test := range tests {
		r, w, _ := os.Pipe()
		log.SetOutput(w)
		defer log.SetOutput(os.Stderr)

		jsonStr := fmt.Sprintf(`{"message": "some string"}`)
		payload := strings.NewReader(jsonStr)

		req := httptest.NewRequest("POST", "/", payload)
		req.Header.Set("Ce-Id", "1234")
		req.Header.Set("Ce-Source", "//storage.googleapis.com/projects/YOUR-PROJECT")
		rr := httptest.NewRecorder()
		GenericHandler(rr, req)

		w.Close()

		if code := rr.Result().StatusCode; code == http.StatusBadRequest {
			t.Errorf("GenericHandler(%q) invalid input, status code (%q)", test.data, code)
		}

		out, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}
		print(out)
		if got := string(out); strings.Contains(got, test.want) != true {
			t.Errorf("\nGenericHandler(%q): \ngot: %q\nwant to contain: %q", test.data, got, test.want)
		}
	}
}
