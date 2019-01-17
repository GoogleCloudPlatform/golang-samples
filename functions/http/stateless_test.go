// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package http

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExecutionCount(t *testing.T) {
	for i := 1; i < 5; i++ {
		req := httptest.NewRequest("GET", "/", strings.NewReader(""))
		rr := httptest.NewRecorder()
		ExecutionCount(rr, req)
		out, err := ioutil.ReadAll(rr.Result().Body)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}
		want := fmt.Sprintf("Instance execution count: %d", i)
		if got := string(out); got != want {
			t.Fatalf("ExecutionCount got %q, want %q", got, want)
		}
	}
}
