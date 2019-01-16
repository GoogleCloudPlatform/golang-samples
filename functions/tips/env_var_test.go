// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package tips

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"
)

func TestEnvVar(t *testing.T) {
	tests := []struct {
		foo string
	}{
		{"bar"},
		{},
	}
	for _, test := range tests {
		os.Setenv("FOO", test.foo)
		rr := httptest.NewRecorder()
		EnvVar(rr, nil)
		out, err := ioutil.ReadAll(rr.Result().Body)
		if err != nil {
			t.Fatalf("EnvVar(%q) error: ioutil.ReadAll: %v", test.foo, err)
		}
		want := fmt.Sprintf("FOO: %q", test.foo)
		if got := string(out); got != want {
			t.Errorf("EnvVar(%s) got %q, want %q", test.foo, got, want)
		}
	}
}
