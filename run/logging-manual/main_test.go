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
	"encoding/json"
	"io/ioutil"
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
			project:     "example",
			traceHeader: "",
			want:        "",
		},
		{
			name:        "no project, trace",
			project:     "",
			traceHeader: "123/456",
			want:        "",
		},
		{
			name:        "project, trace",
			project:     "example",
			traceHeader: "123/456",
			want:        "projects/example/traces/123",
		},
		{
			name:        "project, invalid trace",
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

		r := callHandler(indexHandler, rr, req)

		out, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}

		var e Entry
		if err := json.Unmarshal(out, &e); err != nil {
			t.Errorf("json.Unmarshal: %q", err)
		}

		if e.Trace != test.want {
			t.Errorf("entry(%s): want (%s), got (%s)", test.name, test.want, e.Trace)
		}
	}
}

func callHandler(h func(w http.ResponseWriter, r *http.Request), rr http.ResponseWriter, req *http.Request) *os.File {
	r, w, _ := os.Pipe()
	originalWriter := log.Writer()
	log.SetOutput(w)
	originalFlags := log.Flags()
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	h(rr, req)

	w.Close()
	log.SetOutput(originalWriter)
	log.SetFlags(originalFlags)

	return r
}
