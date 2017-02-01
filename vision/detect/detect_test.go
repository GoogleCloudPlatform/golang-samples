// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestLabel(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	err := detectLabels(&buf, "../testdata/cat.jpg")
	if err != nil {
		t.Fatalf("got %v, want nil err", err)
	}
	got := buf.String()
	if !strings.Contains(got, "cat") {
		t.Errorf("got %q, want to contain cat", got)
	}
}
