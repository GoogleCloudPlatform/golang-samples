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
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	tests := []struct {
		name        string
		project     string
		traceHeader string
		want        string
	}{
		{
			name:        "no project, no trace",
			project:     "",
			traceHeader: "",
			want:        "",
		},
		{
			name:        "no project and trace",
			project:     "",
			traceHeader: "123/456",
			want:        "",
		},
		{
			name:        "project and trace",
			project:     "example",
			traceHeader: "123/456",
			want:        "projects/example/traces/123",
		},
		{
			name:        "project and invalid trace",
			project:     "example",
			traceHeader: "/123",
			want:        "",
		},
	}
	for _, test := range tests {
		projectID = test.project
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("X-Cloud-Trace-Context", test.traceHeader)
		rr := httptest.NewRecorder()

		b := callHandler(indexHandler, rr, req)

		var e Entry
		if err := json.Unmarshal(b.Bytes(), &e); err != nil {
			t.Errorf("json.Unmarshal: %v", err)
		}

		if e.Trace != test.want {
			t.Errorf("indexHandler %q: want %q, got %q", test.name, test.want, e.Trace)
		}
	}
}

// callHandler calls an HTTP handler with the provided request and returns the log output.
func callHandler(h func(w http.ResponseWriter, r *http.Request), rr http.ResponseWriter, req *http.Request) bytes.Buffer {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	originalWriter := os.Stderr
	log.SetOutput(writer)
	defer log.SetOutput(originalWriter)

	h(rr, req)
	writer.Flush()
	return buf
}
