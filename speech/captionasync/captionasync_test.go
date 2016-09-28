// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"golang.org/x/net/context"

	speech "cloud.google.com/go/speech/apiv1beta1"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestRecognize(t *testing.T) {
	testutil.SystemTest(t)

	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	opName, err := send(client, "./quit.raw")
	if err != nil {
		t.Fatal(err)
	}
	if opName == "" {
		t.Fatal("got no op name; want one")
	}
	resp, err := wait(client, opName)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Results) == 0 {
		t.Fatal("got no results; want at least one")
	}
	result := resp.Results[0]
	if len(result.Alternatives) < 1 {
		t.Fatal("got no alternatives; want at least one")
	}
	if got, want := result.Alternatives[0].Transcript, "quit"; got != want {
		t.Errorf("Transcript: got %q; want %q", got, want)
	}
}
