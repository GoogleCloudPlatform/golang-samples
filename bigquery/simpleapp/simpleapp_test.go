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

func TestSimpleApp(t *testing.T) {
	tc := testutil.SystemTest(t)

	rows, err := query(tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	if err := printResults(&b, rows); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(b.String(), "hamlet") {
		t.Errorf("got output: %q; want it to contain hamlet", b.String())
	}
}
