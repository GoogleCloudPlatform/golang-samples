// Copyright 2019 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
