// Copyright 2022 Google LLC
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

package clientendpoint

import (
	"bytes"
	"net/http"
	"os"
	"strings"
	"testing"

	"google.golang.org/api/option"
)

type mockTransport struct {
	gotURL  []string
	gotPath []string
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Mock a success roundtrip to extract request URL
	t.gotURL = append(t.gotURL, req.URL.String())
	t.gotPath = append(t.gotPath, req.URL.Path)

	return &http.Response{StatusCode: 200}, nil
}

func TestSetClientEndpoint(t *testing.T) {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		t.Skip("GOOGLE_APPLICATION_CREDENTIALS not set")
	}

	baseURL := "https://localhost:8080/"
	path := "storage/v1/"
	customEndpoint := baseURL + path

	var buf bytes.Buffer
	mt := mockTransport{}
	opt := option.WithHTTPClient(&http.Client{Transport: &mt})
	if err := setClientEndpoint(&buf, customEndpoint, opt); err != nil {
		t.Errorf("setClientEndpoint: %s", err)
	}
	// Test request URL is set to custom endpoint.
	for _, got := range mt.gotURL {
		if !strings.Contains(got, baseURL) {
			t.Errorf("setClientEndpoint: got request base URL %q; want to contain %q", got, baseURL)
		}
	}
	for _, got := range mt.gotPath {
		if !strings.Contains(got, path) {
			t.Errorf("setClientEndpoint: got request URL path %q; want to contain %q", got, path)
		}
	}
}
