// Copyright 2019 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"image"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	_ "image/png"
)

func TestDiagramHandlerErrors(t *testing.T) {
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
			want:        "Bad Request\n",
			contentType: "text/plain; charset=utf-8",
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/?dot="+url.QueryEscape(test.data), strings.NewReader(""))
		rr := httptest.NewRecorder()
		diagramHandler(rr, req)
		out, err := ioutil.ReadAll(rr.Result().Body)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}

		if got := string(out); got != test.want {
			t.Errorf("png.Decode: response (%s): got %q, want %q", test.label, got, test.want)
		}

		got := rr.Result().Header.Get("Content-Type")
		if got != test.contentType {
			t.Errorf("response (%s) Content-Type: got %q, want %q", test.label, got, test.contentType)
		}
	}
}

func TestDiagramHandlerImage(t *testing.T) {
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
		if (got != test.cacheControl) {
			t.Errorf("response (%s) Cache-Control: got %q, want %q", test.label, got, test.cacheControl)
		}
	}
}
