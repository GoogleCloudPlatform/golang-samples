// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestLabel(t *testing.T) {
	testutil.SystemTest(t)

	labels, err := findLabels("../testdata/cat.jpg")
	if err != nil {
		t.Fatalf("got %v, want nil err", err)
	}
	if len(labels) == 0 {
		t.Fatalf("want non-empty slice of labels")
	}
}
