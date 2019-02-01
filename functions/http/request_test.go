// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMakeRequest(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
	defer s.Close()
	// Override the url global with the test server URL.
	urlString = s.URL

	req := httptest.NewRequest("GET", "/", strings.NewReader(""))
	rr := httptest.NewRecorder()
	MakeRequest(rr, req)
	out, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	want := "ok"
	if got := string(out); got != want {
		t.Errorf("MakeRequest got %q, want %q", got, want)
	}
}
