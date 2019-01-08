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
	err := transcribeMultichannel(&buf, "../testdata/commercial_stereo.wav")
	if err != nil {
		t.Fatal(err)
	}

	if got := buf.String(); !strings.Contains(got, "Channel 1: hi I'd like to buy a Chromecast I'm always wondering whether you could help me with that\nChannel 2: certainly which color would you like we have blue black and red\nChannel 1:  let's go with the black one\nChannel 2:  would you like the new Chromecast Ultra model or the regular Chromecast\nChannel 1:  regular Chromecast is fine thank you\nChannel 2:  okay sure would you like to ship it regular or Express\nChannel 1:  express please\nChannel 2:  terrific it's on the way thank you\nChannel 1:  thank you very much bye\n") {
		t.Fatalf(`transcribeMultichannel(../testdata/commercial_stereo.wav) = %q; want "Channel 1: hi I'd like to buy a Chromecast I'm always wondering whether you could help me with that\nChannel 2: certainly which color would you like we have blue black and red\nChannel 1:  let's go with the black one\nChannel 2:  would you like the new Chromecast Ultra model or the regular Chromecast\nChannel 1:  regular Chromecast is fine thank you\nChannel 2:  okay sure would you like to ship it regular or Express\nChannel 1:  express please\nChannel 2:  terrific it's on the way thank you\nChannel 1:  thank you very much bye\n"`, got)
	}
}
