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

package cloudruntests

import (
	"net/http"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDiagramService(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	service := cloudrunci.NewService("diagram", tc.ProjectID)
	service.Dir = "../system_package"
	if err := service.Deploy(); err != nil {
		t.Fatalf("service.Deploy %q: %v", service.Name, err)
	}
	defer service.Clean()

	requestPath := "/diagram.png"
	req, err := service.NewRequest("GET", requestPath)
	if err != nil {
		t.Fatalf("service.NewRequest: %v", err)
	}
	q := req.URL.Query()
	q.Add("dot", "digraph G { A -> {B, C, D} -> {F} }")
	req.URL.RawQuery = q.Encode()

	resp, err := service.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}

	if got := resp.StatusCode; got != http.StatusOK {
		t.Errorf("response status: got %d, want %d", got, http.StatusOK)
	}

	if got, want := resp.Header.Get("Content-Type"), "image/png"; got != want {
		t.Errorf("response Content-Type: got %q, want %s", got, want)
	}

	if got, want := resp.Header.Get("Cache-Control"), "public, max-age=86400"; got != want {
		t.Errorf("response Cache-Control: got %q, want %q", got, want)
	}
}
