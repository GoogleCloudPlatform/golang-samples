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
package authentication

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func Test_makeGetRequest(t *testing.T) {
	metadata := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.URL.Query().Get("audience")
		if v == "" {
			http.Error(w, "audience is empty", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "TOKEN_FOR_%s", v)
	}))
	defer metadata.Close()

	metadataURL, err := url.Parse(metadata.URL)
	if err != nil {
		t.Fatalf("failed to parse metadata URL (%s): %v", metadata.URL, err)
	}
	defer os.Unsetenv("GCE_METADATA_HOST")
	os.Setenv("GCE_METADATA_HOST", metadataURL.Host)

	var gotHeader string
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("Authorization")
		if gotHeader == "" {
			http.Error(w, "Authorization header is empty", http.StatusBadRequest)
			return
		}
	}))
	defer target.Close()

	url := fmt.Sprintf("%s/foo", target.URL)
	expectedHeader := fmt.Sprintf("Bearer TOKEN_FOR_%s", url)

	resp, err := makeGetRequest(url)
	if err != nil {
		t.Fatalf("makeGetRequest: ")
	}
	if expected := http.StatusOK; resp.StatusCode != expected {
		t.Fatalf("unexpected response status: got=%d; expected=%d", resp.StatusCode, expected)
	}
	if gotHeader != expectedHeader {
		t.Fatalf("unexpected Authorization header: expected=%q; got=%q", expectedHeader, gotHeader)
	}
	defer resp.Body.Close()
}
