// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status: got (%v) want (%v)", status, http.StatusOK)
	}

	want := "Hello, World!"
	if got := rr.Body.String(); !strings.Contains(got, want) {
		t.Errorf("unexpected body: got (%v) want (%v)", got, want)
	}
}

func TestSetup(t *testing.T) {
	setup(context.Background())
	if startupTime.IsZero() {
		t.Error("warmupApp: got (uninitialized startupTime) want (startupTime)")
	}
	if client == nil {
		t.Error("warmupApp: got (uninitialized client) want (client)")
	}
}
