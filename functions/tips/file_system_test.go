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

func TestListFiles(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	ListFiles(rr, req)

	out, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	want := "Files:"
	if got := string(out); !strings.Contains(got, want) {
		t.Errorf("ListFiles got %q, want to contain %q", got, want)
	}
}
