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

	wants := []string{
		"Channel 1: hi I'd like to buy a Chromecast",
		"Channel 2: certainly which color",
	}

	for _, want := range wants {
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf(`transcribeMultichannel(../testdata/commercial_stereo.wav) = \n\n%q\n\nWant to contain \n\n%q`, got, want)
		}
	}
}
