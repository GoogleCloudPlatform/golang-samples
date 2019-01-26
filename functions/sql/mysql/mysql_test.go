// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sql

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMySQLDemo(t *testing.T) {
	if connectionName == "" {
		t.Skip("MySQL database not configured")
	}
	rr := httptest.NewRecorder()
	MySQLDemo(rr, nil)
	out, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll: %v", err)
	}
	want := "Now:"
	if got := string(out); !strings.Contains(got, want) {
		t.Fatalf("MySQLDemo got %q, want to contain %q", got, want)
	}
}
