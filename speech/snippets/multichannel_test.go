// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMultichannel(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	if err := transcribeMultichannel(&buf, "../testdata/commercial_stereo.wav"); err != nil {
		t.Fatal(err)
	}

	var want = "Channel 1: hi I'd like to buy a Chromecast I'm always wondering whether you could help me with that\nChannel 2: certainly which color would you like we have blue black and red"

	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf(`transcribeMultichannel(../testdata/commercial_stereo.wav) = %q; want to contain %q`, got, want)
	}
}
