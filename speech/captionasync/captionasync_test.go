// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"strings"
	"testing"

	"golang.org/x/net/context"

	speech "cloud.google.com/go/speech/apiv1"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestRecognize(t *testing.T) {
	testutil.SystemTest(t)

	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer

	if err := send(&buf, client, "../testdata/quit.raw"); err != nil {
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

	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer

	if err := sendGCS(&buf, client, "gs://python-docs-samples-tests/speech/audio.raw"); err != nil {
		t.Fatal(err)
	}
	if len(buf.String()) == 0 {
		t.Fatal("got no results; want at least one")
	}
	if got, want := buf.String(), "how old is the Brooklyn Bridge"; !strings.Contains(got, want) {
		t.Errorf("Transcript: got %q; want %q", got, want)
	}
}
