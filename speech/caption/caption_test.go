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

func TestRecognize(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer

	if err := recognize(&buf, "../testdata/quit.raw"); err != nil {
		t.Fatal(err)
	}
	if len(buf.String()) == 0 {
		t.Fatal("got no results; want at least one")
	}

	if got, want := buf.String(), "quit"; !strings.Contains(got, want) {
		t.Errorf("Transcript: got %q; want %q", got, want)
	}
}

func TestRecognizeGCS(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer

	if err := recognizeGCS(&buf, "gs://python-docs-samples-tests/speech/audio.raw"); err != nil {
		t.Fatal(err)
	}
	if len(buf.String()) == 0 {
		t.Fatal("got no results; want at least one")
	}

	if got, want := buf.String(), "how old is the Brooklyn Bridge"; !strings.Contains(got, want) {
		t.Errorf("Transcript: got %q; want %q", got, want)
	}
}
