// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package helloworld

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHelloGet(t *testing.T) {
	payload := strings.NewReader("")
	req := httptest.NewRequest("GET", "/", payload)

	rr := httptest.NewRecorder()
	HelloGet(rr, req)

	out, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	want := "Hello, World!"
	if got := string(out); got != want {
		t.Errorf("HelloWorld = %q, want %q", got, want)
	}
}
