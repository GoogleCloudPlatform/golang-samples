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
	"image"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	_ "image/png"
)

func TestDiagramHandlerErrors(t *testing.T) {
	checkGraphviz(t)
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)

	tests := []struct {
		label       string
		data        string
		want        string
		contentType string
	}{
		{
			label:       "empty",
			data:        "",
			want:        "Bad Request\n",
			contentType: "text/plain; charset=utf-8",
		},
		{
			label:       "invalid",
			data:        "digraph",
			want:        "Bad Request: DOT syntax error\n",
			contentType: "text/plain; charset=utf-8",
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/?dot="+url.QueryEscape(test.data), strings.NewReader(""))
		rr := httptest.NewRecorder()
		diagramHandler(rr, req)

		if got := rr.Body.String(); got != test.want {
			t.Errorf("png.Decode: response (%s): got %q, want %q", test.label, got, test.want)
		}

		got := rr.Result().Header.Get("Content-Type")
		if got != test.contentType {
			t.Errorf("response (%s) Content-Type: got %q, want %q", test.label, got, test.contentType)
		}
	}
}

func TestDiagramHandlerImage(t *testing.T) {
	checkGraphviz(t)

	tests := []struct {
		label        string
		data         string
		contentType  string
		cacheControl string
	}{
		{
			label:        "basic diagram",
			data:         "digraph G { A -> {B, C, D} -> {F} }",
			contentType:  "image/png",
			cacheControl: "public, max-age=86400",
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/?dot="+url.QueryEscape(test.data), strings.NewReader(""))
		rr := httptest.NewRecorder()
		diagramHandler(rr, req)
		if _, _, err := image.DecodeConfig(rr.Result().Body); err != nil {
			t.Errorf("image.Decode: response (%s): invalid image: %v", test.label, err)
		}

		got := rr.Result().Header.Get("Content-Type")
		if got != test.contentType {
			t.Errorf("response (%s) Content-Type: got %q, want %s", test.label, got, test.contentType)
		}

		got = rr.Result().Header.Get("Cache-Control")
		if got != test.cacheControl {
			t.Errorf("response (%s) Cache-Control: got %q, want %q", test.label, got, test.cacheControl)
		}
	}
}

func checkGraphviz(t *testing.T) {
	fileInfo, err := os.Stat("/usr/bin/dot")
	if err != nil {
		t.Skipf("os.Stat: %v (install graphviz?)", err)
	}
	if fileInfo.Mode()&0111 == 0 {
		t.Skipf("/usr/bin/dot not executable (install graphviz?)")
	}
}
