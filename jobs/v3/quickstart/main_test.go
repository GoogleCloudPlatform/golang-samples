// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	t.Skip("Need to set GOOGLE_CLOUD_PROJECT")
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = oldStdout

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	got := string(out)

	want := "Request ID"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}
}
