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

func TestSyncLocal(t *testing.T) {
	testutil.SystemTest(t)

	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := syncWords(client, &out, "../testdata/quit.raw"); err != nil {
		t.Fatal(err)
	}
	if got, want := out.String(), `Word: "quit" (startTime=`; !strings.Contains(got, want) {
		t.Errorf("got %q, want to contain %q", got, want)
	}
}

func TestAsyncGCS(t *testing.T) {
	testutil.SystemTest(t)

	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := asyncWords(client, &out, "gs://python-docs-samples-tests/speech/audio.raw"); err != nil {
		t.Fatal(err)
	}
	if got, want := out.String(), `Word: "Brooklyn" (startTime=`; !strings.Contains(got, want) {
		t.Errorf("got %q, want to contain %q", got, want)
	}
	if got, want := out.String(), `Word: "Bridge" (startTime=`; !strings.Contains(got, want) {
		t.Errorf("got %q, want to contain %q", got, want)
	}
}
