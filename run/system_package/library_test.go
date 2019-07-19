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
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	_ "image/png"
)

func TestDiagramLibrary(t *testing.T) {
	checkGraphviz(t)

	files, err := ioutil.ReadDir("library")
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".dot" {
			continue
		}
		f := filepath.Join("library", file.Name())
		out, err := ioutil.ReadFile(f)
		if err != nil {
			t.Errorf("ReadFile (%s): read error: %v", file.Name(), err)
			continue
		}
		req := httptest.NewRequest("GET", "/?dot="+url.QueryEscape(string(out)), strings.NewReader(""))
		rr := httptest.NewRecorder()
		diagramHandler(rr, req)
		if _, _, err := image.DecodeConfig(rr.Result().Body); err != nil {
			t.Errorf("image.DecodeConfig: %s: %v", file.Name(), err)
		}
	}
}
