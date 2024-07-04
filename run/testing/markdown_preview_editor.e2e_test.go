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

package cloudruntests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var (
	editorService *cloudrunci.Service
	renderService *cloudrunci.Service
	client        http.Client
)

func TestEditorService(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	client = http.Client{Timeout: 10 * time.Second}

	renderService = cloudrunci.NewService("renderer", tc.ProjectID)
	renderService.Dir = "../markdown-preview/renderer"
	if err := renderService.Deploy(); err != nil {
		t.Fatalf("service.Deploy %q: %v", renderService.Name, err)
	}
	defer renderService.Clean()

	editorService = cloudrunci.NewService("editor", tc.ProjectID)
	editorService.Dir = "../markdown-preview/editor"
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

	t.Run("UI", caseEditorServiceUI)
	t.Run("Render", caseEditorServiceRender)
}

func caseEditorServiceUI(t *testing.T) {
	req, err := editorService.NewRequest(http.MethodGet, "/")
	if err != nil {
		t.Fatalf("service.NewRequest: %q", err)
	}

	resp, err := editorService.Do(req)
	if err != nil {
		t.Fatalf("client.Do: %v", err)
	}
	defer resp.Body.Close()
	t.Logf("client.Do: %s %s\n", req.Method, req.URL)

	wantStatus := http.StatusOK
	if got := resp.StatusCode; got != wantStatus {
		t.Errorf("response status: got %d, want %d", got, wantStatus)
	}
}

func caseEditorServiceRender(t *testing.T) {
	req, err := editorService.NewRequest(http.MethodPost, "/render")
	if err != nil {
		t.Fatalf("service.NewRequest: %q", err)
	}
	d := struct{ Data string }{Data: "**strong text**"}
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("json.Marshall: %v", err)
	}
	req.Body = io.NopCloser(bytes.NewReader(b))

	resp, err := editorService.Do(req)
	if err != nil {
		t.Fatalf("client.Do: %v", err)
	}
	defer resp.Body.Close()
	t.Logf("client.Do: %s %s\n", req.Method, req.URL)

	wantStatus := http.StatusOK
	if got := resp.StatusCode; got != wantStatus {
		t.Errorf("response status: got %d, want %d", got, wantStatus)
	}

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("io.ReadAll: %v", err)
	}

	want := "<p><strong>strong text</strong></p>\n"
	if got := string(out); got != want {
		t.Errorf("markdown: got %q, want %q", got, want)
	}

}
