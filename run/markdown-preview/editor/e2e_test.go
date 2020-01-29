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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var tests = []struct {
	label      string
	req        *http.Request
	input      string
	want       string
	wantStatus int
}{
	{
		label:      "markdown",
		input:      "**strong text**",
		want:       "<p><strong>strong text</strong></p>\n",
		wantStatus: http.StatusOK,
	},
}

func TestEditorService(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	renderService := cloudrunci.NewService("render", tc.ProjectID)
	renderService.Dir = "../render"
	if err := renderService.Deploy(); err != nil {
		t.Fatalf("service.Deploy %q: %v", renderService.Name, err)
	}
	defer renderService.Clean()

	editorService := cloudrunci.NewService("editor", tc.ProjectID)
	u, err := renderService.URL("")
	if err != nil {
		t.Fatalf("service.URL: %v", err)
	}
	editorService.Env = cloudrunci.EnvVars{
		"EDITOR_UPSTREAM_RENDER_URL": u,
	}
	if err := editorService.Deploy(); err != nil {
		t.Fatalf("service.Deploy %q: %v", editorService.Name, err)
	}
	defer editorService.Clean()

	for _, test := range tests {
		req, err := editorService.NewRequest("POST", "/render")
		if err != nil {
			t.Fatalf("service.NewRequest: %q", err)
		}
		d := struct{ Data string }{Data: test.input}
		b, err := json.Marshal(d)
		if err != nil {
			t.Fatalf("json.Marshall: %v", err)
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(b))

		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("client.Do: %v", err)
		}
		defer resp.Body.Close()
		fmt.Printf("client.Do: %s %s\n", req.Method, req.URL)

		if got := resp.StatusCode; got != test.wantStatus {
			t.Errorf("response status: got %d, want %d", got, test.wantStatus)
		}

		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("ioutil.ReadAll: %v", err)
		}

		if got := string(out); got != test.want {
			t.Errorf("%s: got %q, want %q", test.label, got, test.want)
		}
	}
}
