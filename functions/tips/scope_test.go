// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package tips

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestScopeDemo(t *testing.T) {
	req := httptest.NewRequest("GET", "/", strings.NewReader(""))
	rr := httptest.NewRecorder()
	ScopeDemo(rr, req)
	out, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	want := `Global: "slow", Local: "fast"`
	if got := string(out); got != want {
		t.Errorf("ScopeDemo got %q, want %q", got, want)
	}
}
