// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestReadTimeSeriesAlign(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	if err := readTimeSeriesAlign(buf, tc.ProjectID); err != nil {
		t.Errorf("readTimeSeriesAlign: %v", err)
	}
	want := "Done"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("readTimeSeriesAlign got %q, want to contain %q", got, want)
	}
}

func TestReadTimeSeriesReduce(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	if err := readTimeSeriesReduce(buf, tc.ProjectID); err != nil {
		t.Errorf("readTimeSeriesReduce: %v", err)
	}
	want := "Done"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("readTimeSeriesReduce got %q, want to contain %q", got, want)
	}
}

func TestReadTimeSeriesFields(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	if err := readTimeSeriesFields(buf, tc.ProjectID); err != nil {
		t.Errorf("readTimeSeriesFields: %v", err)
	}
	want := "Done"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("readTimeSeriesFields got %q, want to contain %q", got, want)
	}
}
